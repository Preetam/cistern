var m = require("mithril");
var d3 = require("d3");
var ChartState = require("./ChartState");
var Chart = require("./Chart");
var groupColor = require("./groupColor");

var ChartContainer = {
  oninit: function(vnode) {

    var queryString = m.parseQueryString(window.location.search.replace(/^[?]/, ''))
    var start = new Date(), end = new Date();
    if (queryString.start) {
      start = new Date(queryString.start)
    }
    if (queryString.end) {
      end = new Date(queryString.end)
    }

    var data = {
      series: [],
      query: {
        time_range: {
          start: new Date(),
          end: new Date()
        }
      }
    };
    vnode.state.start = start;
    vnode.state.end = end;
    vnode.state.query = queryString.query || "";

    vnode.state.brushEnd = function(start, end) {
      this.start = start
      this.end = end
      this.refresh()
      this.updateURL()
    }.bind(vnode.state)

    vnode.state.refresh = function() {
      var process = function(data) {
        var columnToName = function(column) {
          return column.aggregate + "(" + column.name + ")";
        };
        var charts = [];
        var chartData = {};
        for (var i in data.query.columns) {
          var column = data.query.columns[i];
          var name = columnToName(column);
          charts.push(name);
          chartData[name] = {};
        }
        for (var i in data.series) {
          var point = data.series[i];
          var groupID = point["_group_id"];
          for (var j in charts) {
            var chartName = charts[j];
            if (!chartData[chartName][groupID]) {
              chartData[chartName][groupID] = [];
            }
            chartData[chartName][groupID].push({
              ts: new Date(point["_ts"]),
              y: point[chartName]
            });
          }
        }
        console.log(JSON.stringify(chartData));

        this.chartStates = {};
        for (var i in chartData) {
          this.chartStates[i] = new ChartState(300, 300, new Date(data.query.time_range.start), new Date(data.query.time_range.end), {
            lines: chartData[i]
          }, i, this.brushEnd);
        }

        this.summaryRows = data.summary;
        this.start = new Date(data.query.time_range.start);
        this.end = new Date(data.query.time_range.end);

        console.log(this.chartStates)
      }.bind(this)

      vnode.state.updateURL = function() {
        var queryString = m.parseQueryString(window.location.search.replace(/^[?]/, ''))
        queryString.start = this.start.toJSON()
        queryString.end = this.end.toJSON()
        queryString.query = this.query

        var newurl = window.location.protocol + "//" + window.location.host + window.location.pathname + '?' + m.buildQueryString(queryString);
        window.history.pushState({path: newurl}, '', newurl);
      }.bind(this)

      m.request({
        method: "POST",
        url: "http://localhost:2020/collections/flowlogs/query?" +
          "start=" + Math.floor(this.start.getTime()/1000) + "&" +
          "end=" + Math.floor(this.end.getTime()/1000) + "&" +
          "query=" + encodeURIComponent(this.query)
      }).then(process);
    }.bind(vnode.state);

    window.onpopstate = (function(e) {
      this.refresh()
    }).bind(this);

    vnode.state.refresh()
  },
  view: function(vnode) {
    var chartComponents = [];
    console.log(vnode.state.chartStates)
    for (var i in vnode.state.chartStates) {
      chartComponents.push(m(new Chart(vnode.state.chartStates[i])));
    }

    var summaryRows = vnode.state.summaryRows;

    if (!summaryRows) {
      summaryRows = [{}];
    }
    var headers = Object.keys(summaryRows[0])
    var summaryTable = m("table.table", [
      m("thead",
        m("tr", headers.map(function(d) {
          if (d == "_group_id") {
            return m("th", "")
          }
          return m("th", d)
        }))
      ),
      m("tbody", summaryRows.map(function(row) {
        return m("tr", Object.keys(row).map(function(k) {
          if (k == "_group_id") {
            return m("td", m("div",
              {
                style: {
                  backgroundColor: groupColor(row[k]),
                  height: "1.5rem",
                  width: "3px",
                  marginRight: "5px"
                }
              },
              ""))
          }
          return m("td", row[k])
        }))
      }))
    ])

    return m("div", {style: "width: 100%;"}, [

      // Start timestamp field
      m("input", {
        onchange: m.withAttr("value", function(v) {
          var d = new Date(v);
          if (!isNaN(d.getTime())) {
            vnode.state.start = new Date(v);
            vnode.state.refresh();
            vnode.state.updateURL();
            return;
            for (var i in vnode.state.chartStates) {
              vnode.state.chartStates[i].start = new Date(v);
            }
          }
        }),
        size: 30,
        value: vnode.state.start.toJSON()
      }),

      // End timestamp field
      m("input", {
        onchange: m.withAttr("value", function(v) {
          var d = new Date(v);
          if (!isNaN(d.getTime())) {
            vnode.state.end = new Date(v);
            vnode.state.refresh();
            vnode.state.updateURL();
            return;
            for (var i in vnode.state.chartStates) {
              vnode.state.chartStates[i].end = new Date(v);
            }
          }
        }),
        size: 30,
        value: vnode.state.end.toJSON()
      }),

      // Query field
      m("input", {
        style: {display: "block"},
        onchange: m.withAttr("value", function(v) {
          vnode.state.query = v;
          vnode.state.refresh();
          vnode.state.updateURL();
        }),
        size: 120,
        value: vnode.state.query
      }),

      // Charts
      chartComponents,

      // Summary table
      summaryTable
    ]);
  }
};

module.exports = ChartContainer;
