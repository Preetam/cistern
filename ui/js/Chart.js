var m = require("mithril");
var d3 = require("d3");
var groupColor = require("./groupColor");

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
        var color = groupColor(i);
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

module.exports = Chart;
