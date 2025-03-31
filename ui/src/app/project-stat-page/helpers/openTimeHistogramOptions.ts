import { Options } from 'highcharts';

export const openTimeHistogramOptions: Options = {
  chart: {
    type: 'column',
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Time in Open State',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Days'
    }
  },
  legend: {
    enabled: false,
  },
  xAxis: {
    lineColor: '#fff',
    categories: [],
    title: {
      text: 'Time intervals'
    }
  },
  plotOptions: {
    series: {
      borderRadius: 5,
    } as any,
  },
  series: [
    {
      type: 'column',
      color: '#506ef9',
      data: [],
    },
  ],
};
