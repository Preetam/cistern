var m = require("mithril");

var d3 = require("d3");

var App = {
  view: function(vnode) {
    return m("div", "todo");
  }
};

var groupIDToColorString = function(groupID) {
  var hue = parseInt(groupID, 16)
  hue = hue % 360
  return "hsl(" + hue + ", 50%, 50%)"
}

var ChartState = function(width, height, start, end, data, name, brushEnd) {
  this.width = width;
  this.height = height;
  this.start = start;
  this.end = end;
  this.data = data;
  this.name = name;
  this.maxVal = 10;
  this.brushEnd = brushEnd
  for (var i in data.lines) {
    var lineData = data.lines[i];
    if (lineData.length == 1) {
      // Only one point, so skip it.
      continue
    }
    for (var j in lineData) {
      if (lineData[j].y > this.maxVal) {
        console.log("max val", name, i, lineData[j].y)
        this.maxVal = lineData[j].y;
      }
    }
  }
};

var Chart = function(chartState) {
  this.oninit = function(vnode) {
    vnode.state.chartState = chartState;
  };
  this.view = function(vnode) {
    // Resize
    var resize = function(vnode) {
      var chart = vnode.dom;
      var div = d3.select(chart.parentNode).select(".tooltip");
      var width = parseInt(d3.select(chart).style("width"));
      var data = this.chartState.data, w = width, h = this.chartState.height, margin = 70, y = d3.scaleLinear().domain([ 0, this.chartState.maxVal * 1.1 ]).range([ h - margin, 0 + margin ]), x = d3.scaleTime().domain([ this.chartState.start, this.chartState.end ]).range([ 0 + margin, w - margin ]);
      var yAxis = d3.axisLeft(y).ticks(3).tickFormat(d3.format(".0s"));
      var xAxis = d3.axisBottom(x).ticks(4);
      // Remove existing paths
      d3.select(chart).selectAll("path").remove();
      // Draw paths
      for (i in data.lines) {
        var lineData = data.lines[i];
        var line = d3.line().x(function(d, i) {
          return x(d.ts);
        }).y(function(d) {
          return y(d.y);
        });
        var color = groupIDToColorString(i);
        d3.select(chart).select(".lineGroup").append("path").attr("d", line(lineData)).attr("fill", "none").attr("stroke", color).attr("stroke-width", "1px");
      }
      // Draw axes
      d3.select(chart).select(".y-axis").attr("transform", "translate(" + (margin - 20) + ", 0)").call(yAxis);
      d3.select(chart).select(".x-axis").attr("transform", "translate(0, " + (h - margin + 10) + ")").call(xAxis);

      // Set up brush
      brushended = function() {
        var s = d3.event.selection;
        if (s) {
            var start = x.invert(s[0]);
            var end = x.invert(s[1]);
            this.chartState.brushEnd(start, end)
        } else {
            var end = new Date();
            var start = new Date(end - 90*86400*1000);
            this.chartState.brushEnd(start, end)
        }
      }.bind(this);
      var brush = d3.brushX().on("end", brushended).extent([ [ margin, margin ], [ w - margin, h - margin ] ]), idleTimeout, idleDelay = 350;
      d3.select(chart).select(".brush").call(brush);
    }.bind(this);
    // Draw
    var draw = function(vnode) {
      console.log("this.chartState.name = ", this.chartState.name);
      d3.select(window).on("resize." + this.chartState.name, resize.bind(null, vnode));
      resize(vnode);
    }.bind(this);
    // Elements
    return [
      m("h4", this.chartState.name),
      m("svg", {
        width: "100%",
        height: this.chartState.height,
        oncreate: draw.bind(this)
      },
      m("g", [ m("g.lineGroup"), m("g.x-axis"), m("g.y-axis"), m("g.brush") ]))
    ];
  };
};

var ChartContainer = {
  oninit: function(vnode) {
    var data = {
      series: [],
      query: {
        time_range: {
          start: new Date(),
          end: new Date()
        }
      }
    };
    vnode.state.start = data.query.time_range.start;
    vnode.state.end = data.query.time_range.end;

    vnode.state.brushEnd = function(start, end) {
      this.start = start
      this.end = end
      this.refresh()
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

      m.request({
        method: "POST",
        url: "http://localhost:2020/collections/flowlogs/query?" +
          "start=" + Math.floor(this.start.getTime()/1000) + "&" +
          "end=" + Math.floor(this.end.getTime()/1000) + "&" +
          "query=" + encodeURIComponent(this.query)
      }).then(process);
    }.bind(vnode.state);
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
                  backgroundColor: groupIDToColorString(row[k]),
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

var CollectionChartPage = {
  view: function(vnode) {
    return m("div", [ m("div", m(ChartContainer)) ]);
  }
};

m.mount(document.getElementById("app"), CollectionChartPage);
