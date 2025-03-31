import { Options } from 'highcharts';

export const priorityDistributionOptions: Options = {
  chart: {
    type: 'bar',
  },
  credits: {
    enabled: false,
  },
  title: {
    text: 'Task Priority Distribution',
  },
  yAxis: {
    visible: true,
    title: {
      text: 'Issue count'
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
      type: 'bar',
      color: '#9C27B0',
      data: [],
    },
  ],
};
