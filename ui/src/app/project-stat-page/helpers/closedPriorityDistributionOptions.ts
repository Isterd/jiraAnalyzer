import { Options } from 'highcharts';

export const closedPriorityDistributionOptions: Options = {
  chart: {
    type: 'column',
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Closed Tasks Priority',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Closed issue count'
    }
  },
  legend: {
    enabled: false,
  },
  xAxis: {
    lineColor: '#fff',
    categories: [],
    title: {
      text: 'Priority'
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
      color: '#E91E63',
      data: [],
    },
  ],
};
