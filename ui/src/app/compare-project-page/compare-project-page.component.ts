import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DatabaseProjectServices } from '../services/database-project.services';

interface ProjectMetrics {
  total_issues?: number;
  open_issues?: number;
  closed_issues?: number;
  reopen_issues?: number;
  resolved_issues?: number;
  in_progress_issues?: number;
  average_time_issues?: number;
  average_count_issues?: number;

  // Альтернативные названия полей
  allIssuesCount?: number;
  openIssuesCount?: number;
  closeIssuesCount?: number;
  reopenedIssuesCount?: number;
  resolvedIssuesCount?: number;
  progressIssuesCount?: number;
  averageTime?: number;
  averageIssuesCount?: number;
}

interface ProjectComparison {
  key: string;
  metrics: ProjectMetrics;
}

@Component({
  selector: 'app-compare-project-page',
  templateUrl: './compare-project-page.component.html',
  styleUrls: ['./compare-project-page.component.css']
})
export class CompareProjectPageComponent implements OnInit {
  projectKeys: string[] = [];
  comparisonData: ProjectComparison[] = [];
  isLoading = true;
  error: string | null = null;

  metrics = [
    { id: 1, name: 'Общее количество задач', key: 'total_issues', altKey: 'allIssuesCount' },
    { id: 2, name: 'Открытые задачи', key: 'open_issues', altKey: 'openIssuesCount' },
    { id: 3, name: 'Закрытые задачи', key: 'closed_issues', altKey: 'closeIssuesCount' },
    { id: 4, name: 'Переоткрытые задачи', key: 'reopen_issues', altKey: 'reopenedIssuesCount' },
    { id: 5, name: 'Задачи в работе', key: 'in_progress_issues', altKey: 'progressIssuesCount' },
    {
      id: 6,
      name: 'Среднее время (часы)',
      key: 'average_time_issues',
      altKey: 'averageTime',
      isNumber: true
    },
    {
      id: 7,
      name: 'Среднее количество в день',
      key: 'average_count_issues',
      altKey: 'averageIssuesCount',
      isNumber: true
    }
  ];

  constructor(
    private route: ActivatedRoute,
    private projectService: DatabaseProjectServices
  ) {}

  ngOnInit(): void {
    this.route.queryParamMap.subscribe(params => {
      this.projectKeys = params.get('keys')?.split(',') || [];

      if (this.projectKeys.length >= 2 && this.projectKeys.length <= 3) {
        this.loadComparisonData();
      } else {
        this.error = 'Неверное количество проектов для сравнения (требуется 2 или 3)';
        this.isLoading = false;
      }
    });
  }

  loadComparisonData(): void {
    this.projectKeys.forEach(key => {
      this.projectService.getProjectAnalytics(key).subscribe({
        next: (response: any) => {
          console.log('Данные проекта', key, response);
          this.comparisonData.push({
            key,
            metrics: response.data
          });

          if (this.comparisonData.length === this.projectKeys.length) {
            this.isLoading = false;
          }
        },
        error: (err) => {
          console.error('Ошибка загрузки данных для проекта', key, err);
          this.error = `Ошибка загрузки данных для проекта ${key}`;
          this.isLoading = false;
        }
      });
    });
  }

  getMetricValue(project: ProjectComparison, metricKey: string, altKey?: string): any {
    const metrics = project.metrics || {};

    // Проверяем основное ключевое поле
    if (metrics.hasOwnProperty(metricKey)) {
      return metrics[metricKey as keyof ProjectMetrics];
    }

    // Проверяем альтернативное ключевое поле
    if (altKey && metrics.hasOwnProperty(altKey)) {
      return metrics[altKey as keyof ProjectMetrics];
    }

    return 'N/A';
  }

  isNumber(value: any): boolean {
    return typeof value === 'number';
  }
}
