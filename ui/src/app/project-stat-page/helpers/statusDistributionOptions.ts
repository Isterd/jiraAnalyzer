import {Options} from "highcharts";

export const statusDistributionOptions: Options = {
  chart: {
    type: 'pie',
    height: 400
  },
  title: {
    text: 'Распределение задач по состояниям'
  },
  plotOptions: {
    pie: {
      allowPointSelect: true,
      cursor: 'pointer',
      dataLabels: {
        enabled: true,
        format: '<b>{point.name}</b>: {point.y} задач'
      },
      showInLegend: true
    }
  },
  tooltip: {
    pointFormat: '{series.name}: <b>{point.y}</b>'
  },
  accessibility: {
    point: {
      valueSuffix: ' задач'
    }
  }
};
