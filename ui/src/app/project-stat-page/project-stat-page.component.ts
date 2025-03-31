import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DatabaseProjectServices, ConfigurationService } from '../services';
import { Chart } from 'angular-highcharts';

import {
  openTimeHistogramOptions,
  statusDistributionOptions,
  activityGraphOptions,
  complexityGraphOptions,
  priorityDistributionOptions,
  closedPriorityDistributionOptions
} from './helpers';

@Component({
  selector: 'app-project-stat-page',
  templateUrl: './project-stat-page.component.html',
  styleUrls: ['./project-stat-page.component.css'],
})
export class ProjectStatPageComponent implements OnInit {
  projects: string[] = [];
  ids: number[] = [];
  webUrl: string = '';

  // Только необходимые графики (6 штук)
  openTimeHistogram: Chart | undefined;
  statusDistributionChart: Chart | undefined;
  activityGraph: Chart | undefined;
  complexityGraph: Chart | undefined;
  priorityDistributionChart: Chart | undefined;
  closedPriorityDistributionChart: Chart | undefined;

  constructor(
    private route: ActivatedRoute,
    private dbProjectService: DatabaseProjectServices,
    private configurationService: ConfigurationService
  ) {}

  resetCharts(): void {
    // Очистка опций графиков
    openTimeHistogramOptions.series = [];
    statusDistributionOptions.series = [];
    activityGraphOptions.series = [];
    complexityGraphOptions.series = [];
    priorityDistributionOptions.series = [];
    closedPriorityDistributionOptions.series = [];

    // Очистка элементов DOM
    const elements = [
      'open-time-histogram',
      'status-distribution',
      'activity-graph',
      'complexity-graph',
      'priority-distribution',
      'closed-priority-distribution',
    ];
    elements.forEach((id) => {
      const elem = document.getElementById(id);
      if (elem) elem.innerHTML = '';
    });
  }

  ngOnInit(): void {
    this.webUrl = this.configurationService.getValue('webUrl');

    // Инициализация только нужных графиков
    this.openTimeHistogram = new Chart(openTimeHistogramOptions);
    this.statusDistributionChart = new Chart(statusDistributionOptions);
    this.activityGraph = new Chart(activityGraphOptions);
    this.complexityGraph = new Chart(complexityGraphOptions);
    this.priorityDistributionChart = new Chart(priorityDistributionOptions);
    this.closedPriorityDistributionChart = new Chart(closedPriorityDistributionOptions);

    this.route.paramMap.subscribe((params) => {
      const projectId = params.get('id');
      if (!projectId) {
        alert('ID проекта не указан.');
        return;
      }
      this.projects = [projectId];
    });

    this.route.queryParamMap.subscribe((queryParams) => {
      const selectedTasks = queryParams.get('value');
      if (!selectedTasks) {
        alert('Недостаточно данных для загрузки графиков.');
        return;
      }
      this.ids = selectedTasks.split(',').map(Number);
      this.loadGraphData();
    });
  }

  loadGraphData(): void {
    if (!this.projects || this.projects.length === 0) return;

    const projectKey = this.projects[0];

    this.ids.forEach(taskNumber => {
      this.dbProjectService.getGraph(taskNumber.toString(), projectKey).subscribe(
        (response: any) => {
          if (response?.data) {
            this.updateChart(taskNumber, response.data);
          }
        },
        (error) => console.error('Ошибка загрузки графика:', error)
      );
    });
  }

  updateChart(taskNumber: number, data: any): void {
    switch (taskNumber) {
      case 1:
        this.updateOpenTimeHistogram(data);
        break;
      case 2:
        this.updateStatusDistribution(data);
        break;
      case 3:
        this.updateActivityGraph(data);
        break;
      case 4:
        this.updateComplexityGraph(data);
        break;
      case 5:
        this.updatePriorityDistribution(data);
        break;
      case 6:
        this.updateClosedPriorityDistribution(data);
        break;
      default:
        console.error(`Неизвестная задача: ${taskNumber}`);
    }
  }

  updateOpenTimeHistogram(data: any): void {
    const container = document.getElementById('open-time-histogram');
    const title = document.getElementById('open-time-histogram-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) {
        title.textContent = 'Гистограмма времени в открытом состоянии - нет данных';
      }
      return;
    }

    const categories = data.map((item: any) => item.DayInterval);
    const values = data.map((item: any) => item.TaskCount);

    // Проверка существования xAxis
    if (!openTimeHistogramOptions.xAxis) {
      console.error('xAxis is undefined');
      return;
    }

    if (Array.isArray(openTimeHistogramOptions.xAxis)) {
      openTimeHistogramOptions.xAxis[0].categories = categories;
    } else {
      openTimeHistogramOptions.xAxis.categories = categories;
    }

    openTimeHistogramOptions.series = [{
      name: this.projects[0],
      type: 'column',
      data: values
    }];

    if (this.openTimeHistogram) {
      this.openTimeHistogram = new Chart({...openTimeHistogramOptions});
    }
  }

  updateStatusDistribution(data: any): void {
    const container = document.getElementById('status-distribution');
    const title = document.getElementById('status-distribution-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) title.textContent = 'Распределение по состояниям - нет данных';
      return;
    }

    // Группируем данные по статусам
    const statusMap = new Map<string, number>();
    data.forEach((item: any) => {
      const current = statusMap.get(item.Status) || 0;
      statusMap.set(item.Status, current + item.TaskCount);
    });

    const seriesData = Array.from(statusMap.entries()).map(([status, count]) => ({
      name: status,
      y: count
    }));

    statusDistributionOptions.series = [{
      type: 'pie',
      name: 'Tasks by Status',
      data: seriesData
    }];

    if (this.statusDistributionChart) {
      this.statusDistributionChart = new Chart({...statusDistributionOptions});
    }
  }

  updateActivityGraph(data: any): void {
    const container = document.getElementById('activity-graph');
    const title = document.getElementById('activity-graph-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) title.textContent = 'График активности - нет данных';
      return;
    }

    // Форматируем даты для лучшей читаемости
    const categories = data.map((item: any) =>
      new Date(item.Day).toLocaleDateString('ru-RU')
    );

    const openedData = data.map((item: any) => item.CumulativeOpened);
    const closedData = data.map((item: any) => item.CumulativeClosed);

    // Проверка существования xAxis
    if (!activityGraphOptions.xAxis) {
      console.error('xAxis is undefined');
      return;
    }

    // Обновляем опции графика
    if (Array.isArray(activityGraphOptions.xAxis)) {
      activityGraphOptions.xAxis[0].categories = categories;
    } else {
      activityGraphOptions.xAxis.categories = categories;
    }

    activityGraphOptions.series = [
      {
        name: 'Открытые',
        type: 'spline',
        data: openedData,
        color: '#4CAF50'
      },
      {
        name: 'Закрытые',
        type: 'spline',
        data: closedData,
        color: '#F44336'
      }
    ];

    if (this.activityGraph) {
      this.activityGraph = new Chart({...activityGraphOptions});
    }
  }

  updateComplexityGraph(data: any): void {
    const container = document.getElementById('complexity-graph');
    const title = document.getElementById('complexity-graph-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) {
        title.textContent = 'Сложность задач - нет данных';
      }
      return;
    }

    const categories = data.map((item: any) => item.ComplexityLevel);
    const values = data.map((item: any) => item.TaskCount);

    // Проверка существования xAxis
    if (!complexityGraphOptions.xAxis) {
      console.error('xAxis is undefined');
      return;
    }

    if (Array.isArray(complexityGraphOptions.xAxis)) {
      complexityGraphOptions.xAxis[0].categories = categories;
    } else {
      complexityGraphOptions.xAxis.categories = categories;
    }

    complexityGraphOptions.series = [{
      name: 'Количество задач',
      type: 'column',
      data: values
    }];

    if (this.complexityGraph) {
      this.complexityGraph = new Chart({...complexityGraphOptions});
    }
  }

  updatePriorityDistribution(data: any): void {
    const container = document.getElementById('priority-distribution');
    const title = document.getElementById('priority-distribution-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) {
        title.textContent = 'Приоритетность задач - нет данных';
      }
      return;
    }

    const categories = data.map((item: any) => item.Priority);
    const values = data.map((item: any) => item.TaskCount);

    // Проверка существования xAxis
    if (!priorityDistributionOptions.xAxis) {
      console.error('xAxis is undefined');
      return;
    }

    if (Array.isArray(priorityDistributionOptions.xAxis)) {
      priorityDistributionOptions.xAxis[0].categories = categories;
    } else {
      priorityDistributionOptions.xAxis.categories = categories;
    }

    priorityDistributionOptions.series = [{
      name: 'Количество задач',
      type: 'bar',
      data: values
    }];

    if (this.priorityDistributionChart) {
      this.priorityDistributionChart = new Chart({...priorityDistributionOptions});
    }
  }

  updateClosedPriorityDistribution(data: any): void {
    const container = document.getElementById('closed-priority-distribution');
    const title = document.getElementById('closed-priority-distribution-title');

    if (!data || data.length === 0) {
      if (container) container.remove();
      if (title) {
        title.textContent = 'Приоритетность закрытых задач - нет данных';
      }
      return;
    }

    const categories = data.map((item: any) => item.Priority);
    const values = data.map((item: any) => item.TaskCount);

    // Проверка существования xAxis
    if (!closedPriorityDistributionOptions.xAxis) {
      console.error('xAxis is undefined');
      return;
    }

    if (Array.isArray(closedPriorityDistributionOptions.xAxis)) {
      closedPriorityDistributionOptions.xAxis[0].categories = categories;
    } else {
      closedPriorityDistributionOptions.xAxis.categories = categories;
    }

    closedPriorityDistributionOptions.series = [{
      name: 'Закрытые задачи',
      type: 'column',
      data: values
    }];

    if (this.closedPriorityDistributionChart) {
      this.closedPriorityDistributionChart = new Chart({...closedPriorityDistributionOptions});
    }
  }
}
