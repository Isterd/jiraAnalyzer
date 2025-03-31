import { Options } from 'highcharts';

export const activityGraphOptions: Options = {
  chart: {
    type: 'spline',
    height: 500
  },
  title: {
    text: 'Активность по задачам'
  },
  xAxis: {
    type: 'category',
    labels: {
      rotation: -45,
      style: {
        fontSize: '10px'
      }
    }
  },
  yAxis: {
    title: {
      text: 'Количество задач'
    },
    min: 0
  },
  tooltip: {
    headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
    pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
      '<td style="padding:0"><b>{point.y}</b></td></tr>',
    footerFormat: '</table>',
    shared: true,
    useHTML: true
  },
  plotOptions: {
    spline: {
      marker: {
        radius: 4,
        lineColor: '#666666',
        lineWidth: 1
      }
    }
  }
};
