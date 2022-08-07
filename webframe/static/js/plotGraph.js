window.onload = function () {

var dataPoints1 = [];
var dataPoints2 = [];

var chart = new CanvasJS.Chart("chartContainer", {
	zoomEnabled: true,
	title: {
		text: "Download and Upload speeds"
	},
	axisX: {
		title: "chart updates every 5 minutes"
	},
	axisY:{
		prefix: "Mbps ",
		includeZero: true
	},
	toolTip: {
		shared: true
	},
	legend: {
		cursor:"pointer",
		verticalAlign: "top",
		fontSize: 22,
		fontColor: "dimGrey",
		itemclick : toggleDataSeries
	},
	data: [{
		type: "line",
		xValueType: "dateTime",
		yValueFormatString: "####Mbps",
		xValueFormatString: "hh:mm:ss TT",
		showInLegend: true,
		name: "Download",
		dataPoints: dataPoints1
		},
		{
			type: "line",
			xValueType: "dateTime",
			yValueFormatString: "####Mbps",
			showInLegend: true,
			name: "Upload" ,
			dataPoints: dataPoints2
	}]
});

function toggleDataSeries(e) {
	if (typeof(e.dataSeries.visible) === "undefined" || e.dataSeries.visible) {
		e.dataSeries.visible = false;
	}
	else {
		e.dataSeries.visible = true;
	}
	chart.render();
}

var updateInterval = 300000;
// initial value
var yValue1 = 14;
var yValue2 = 16;
var xValue = new Date;
var time = new Date;
// starting at 7.30 am
time.setHours(7);
time.setMinutes(30);
time.setSeconds(00);
time.setMilliseconds(00);

function updateChart() {
	{% for row in speedtests %} 
		xValue = Date.parse("{{ row[1] }}");
		yValue1 = {{ row[5] }};
		yValue2 = {{ row[4] }};

	// pushing the new values
	dataPoints1.push({
		x: xValue,
		y: yValue1
	});

	dataPoints2.push({
		x: xValue,
		y: yValue2
	});	
                {% endfor  %}

	// updating legend text with  updated with y Value
	chart.options.data[0].legendText = " Download " + yValue1 + " Mbps";
	chart.options.data[1].legendText = " Upload " + yValue2 + " Mbps";
	chart.render();
}
// generates first set of dataPoints
updateChart();
setInterval(function(){updateChart()}, updateInterval);

}
