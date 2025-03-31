import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IProj } from "../../models/proj.model";
import { SettingBox } from "../../models/setting.model";
import { Router } from "@angular/router";
import { DatabaseProjectServices } from "../../services/database-project.services";
import {catchError, firstValueFrom} from "rxjs";
import {CheckedSetting} from "../../models/check-setting.model";


@Component({
  selector: 'app-my-project',
  templateUrl: './my-project.component.html',
  styleUrls: ['./my-project.component.css']
})
export class MyProjectComponent implements OnInit {
  @Output() onChecked = new EventEmitter<any>();
  @Input() myProject!: IProj;
  stat: ProjectStat = new ProjectStat();
  checked = 0;
  completed = 0;
  processed = false;
  settings = false;
  checkboxes: SettingBox[] = [];
  loadingStats = false;
  loadingAnalysis = false;

  constructor(
    private router: Router,
    private dbProjectService: DatabaseProjectServices
  ) {}

  ngOnInit(): void {
    this.initializeCheckboxes();
    this.loadProjectStatistics();
  }

  private initializeCheckboxes(): void {
    this.checkboxes = [
      new SettingBox("Гистограмма времени в открытом состоянии", false, 1),
      new SettingBox("Диаграммы распределения по состояниям", false, 2),
      new SettingBox("График активности по задачам", false, 3),
      new SettingBox("График сложности задач", false, 4),
      new SettingBox("Приоритетность всех задач", false, 5),
      new SettingBox("Приоритетность закрытых задач", false, 6)
    ];
  }

  private loadProjectStatistics(): void {
    this.loadingStats = true;
    console.log('Загрузка аналитики для проекта:', this.myProject.key);

    this.dbProjectService.getProjectAnalytics(this.myProject.key).subscribe({
      next: (response) => {
        console.log('Обработанные данные аналитики:', response);
        if (response?.data) {
          this.updateStatistics(response.data);
        } else {
          console.warn('Нет данных аналитики для проекта', this.myProject.key);
          // Установите значения по умолчанию
          this.updateStatistics({
            allIssuesCount: 0,
            openIssuesCount: 0,
            closeIssuesCount: 0,
            reopenedIssuesCount: 0,
            resolvedIssuesCount: 0,
            progressIssuesCount: 0,
            averageTime: 0,
            averageIssuesCount: 0
          });
        }
        this.loadingStats = false;
      },
      error: (error) => {
        console.error('Ошибка загрузки аналитики:', error);
        this.loadingStats = false;
      }
    });
  }

  private updateStatistics(data: any): void {
    console.log('Обновление статистики с данными:', data);
    this.stat = {
      AllIssuesCount: data.allIssuesCount ?? 0,
      OpenIssuesCount: data.openIssuesCount ?? 0,
      CloseIssuesCount: data.closeIssuesCount ?? 0,
      ReopenedIssuesCount: data.reopenedIssuesCount ?? 0,
      ResolvedIssuesCount: data.resolvedIssuesCount ?? 0,
      ProgressIssuesCount: data.progressIssuesCount ?? 0,
      AverageTime: data.averageTime ?? 0,
      AverageIssuesCount: data.averageIssuesCount ?? 0
    };

    // Для отладки - временный вывод
    console.log('Обновленная статистика:', this.stat);
  }

  async processProject(): Promise<void> {
    const selectedTasks = this.getSelectedTasks();

    if (selectedTasks.length === 0) {
      alert("Выберите хотя бы одну аналитическую задачу");
      return;
    }

    this.loadingAnalysis = true;
    this.processed = false; // Используем processed вместо tasksProcessed

    try {
      for (const taskNumber of selectedTasks) {
        const response = await firstValueFrom(
          this.dbProjectService.makeGraph(taskNumber.toString(), this.myProject.key)
        );

        if (!response.success) {
          throw new Error(response.message || 'Ошибка сервера');
        }
      }

      this.processed = true; // Устанавливаем processed
      alert("Все выбранные аналитические задачи успешно обработаны!");
    } catch (error) {
      console.error('Ошибка обработки задач:', error);
      alert('Произошла ошибка при обработке: ' + (error instanceof Error ? error.message : 'Неизвестная ошибка'));
    } finally {
      this.loadingAnalysis = false;
    }
  }

  private getSelectedTasks(): number[] {
    return this.checkboxes
      .filter(box => box.Checked)
      .map(box => Number(box.BoxId));
  }

  private prepareForProcessing(): void {
    this.loadingAnalysis = true;
    this.processed = false;
    this.completed = 0;
  }


  private handleProcessingError(error: any): void {
    let errorMessage = 'Произошла ошибка при обработке задач';

    if (error instanceof Error) {
      errorMessage = error.message;
    } else if (error?.message) {
      errorMessage = error.message;
    }

    console.error('Ошибка обработки:', error);
    alert(errorMessage);
    this.resetProcessingState();
  }

  private resetProcessingState(): void {
    this.processed = false;
    this.completed = 0;
    this.loadingAnalysis = false;
    this.checked = 0; // Сбрасываем счетчик выбранных задач
    // Сбрасываем чекбоксы
    this.checkboxes.forEach(box => box.Checked = false);
  }

  checkResult(): void {
    if (!this.processed) { // Проверяем processed
      alert('Сначала завершите обработку выбранных задач');
      return;
    }

    const selectedIds = this.checkboxes
      .filter(box => box.Checked)
      .map(box => box.BoxId);

    if (selectedIds.length === 0) {
      alert('Нет выбранных задач для просмотра');
      return;
    }

    this.navigateToProjectStats(selectedIds);
  }



  private navigateToProjectStats(ids: number[]): void {
    this.router.navigate(['/project-stat', this.myProject.key], {
      queryParams: {
        value: ids.join(',')
      }
    });
  }

  toggleSettings(): void {
    this.settings = !this.settings;
  }

  get canViewResults(): boolean {
    return this.processed && this.checked > 0;
  }

  get canAnalyze(): boolean {
    return this.checked > 0 && !this.loadingAnalysis;
  }

  get hasSelectedTasks(): boolean {
    return this.checked > 0;
  }

  onSettingChanged(setting: CheckedSetting): void {
    const box = this.checkboxes.find(b => b.BoxId === setting.BoxId);
    if (box) {
      box.Checked = setting.Checked;
      this.checked = this.checkboxes.filter(b => b.Checked).length;
    }
  }

}

class ProjectStat {
  AllIssuesCount = 0;
  AverageIssuesCount = 0;
  AverageTime = 0;
  CloseIssuesCount = 0;
  OpenIssuesCount = 0;
  ResolvedIssuesCount = 0;
  ReopenedIssuesCount = 0;
  ProgressIssuesCount = 0;
}
