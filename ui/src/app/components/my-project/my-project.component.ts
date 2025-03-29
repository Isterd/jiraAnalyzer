import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IProj } from "../../models/proj.model";
import { SettingBox } from "../../models/setting.model";
import { Router } from "@angular/router";
import { DatabaseProjectServices } from "../../services/database-project.services";

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
  complited = 0;
  processed = false;
  settings = false;
  checkboxes: SettingBox[] = [];

  constructor(private router: Router, private dbProjectService: DatabaseProjectServices) {}

  ngOnInit(): void {
    this.processed = false;
    this.settings = false;

    // Инициализация чекбоксов
    this.checkboxes.push(new SettingBox("Гистограмма, отражающая время, которое задачи провели в открытом состоянии", false, 1));
    this.checkboxes.push(new SettingBox("Диаграммы, которые показывают распределение времени по состояниям задач", false, 2));
    this.checkboxes.push(new SettingBox("График активности по задачам", false, 3));
    this.checkboxes.push(new SettingBox("График сложности задач", false, 4));
    this.checkboxes.push(new SettingBox("График, отражающий приоритетность всех задач", false, 5));
    this.checkboxes.push(new SettingBox("График, отражающий приоритетность закрытых задач", false, 6));

    // Получение аналитики проекта
    this.dbProjectService.getProjectAnalytics(this.myProject.key).subscribe(
      (response: any) => {
        if (response && response.data) {
          const data = response.data;
          this.stat.AllIssuesCount = data.total_issues || 0;
          this.stat.OpenIssuesCount = data.open_issues || 0;
          this.stat.CloseIssuesCount = data.closed_issues || 0;
          this.stat.ReopenedIssuesCount = data.reopen_issues || 0;
          this.stat.ResolvedIssuesCount = data.resolved_issues || 0;
          this.stat.ProgressIssuesCount = data.in_progress_issues || 0;
          this.stat.AverageTime = data.average_time_issues || 0;
          this.stat.AverageIssuesCount = data.average_count_issues || 0;
        } else {
          console.error("Ответ сервера не содержит данных:", response);
          alert("Не удалось получить данные проекта. Попробуйте позже.");
        }
      },
      (error: any) => {
        console.error("Ошибка при получении статистики проекта:", error);
        alert("Произошла ошибка при загрузке данных проекта.");
      }
    );
  }

  async processProject() {
    const selectedTasks = this.checkboxes
      .filter((box) => box.Checked)
      .map((box) => Number(box.BoxId));

    if (selectedTasks.length === 0) {
      alert("Выберите хотя бы одну аналитическую задачу.");
      return;
    }

    this.processed = true;
    this.complited = 0;

    try {
      for (const taskNumber of selectedTasks) {
        await this.dbProjectService.makeGraph(taskNumber.toString(), this.myProject.key).toPromise();
        this.complited++;
      }
      alert("Все графики успешно созданы.");
    } catch (error) {
      console.error("Ошибка при создании графиков:", error);
      alert("Произошла ошибка при создании графиков.");
    }
  }

  checkResult() {
    const ids = this.checkboxes
      .filter((box) => box.Checked)
      .map((box) => Number(box.BoxId));

    this.router.navigate(['/project-stat'], {
      queryParams: {
        keys: this.myProject.key,
        value: ids.join(','),
      },
    });
  }

  clickOnSettings() {
    this.settings = !this.settings;
  }

  disableCheckResultButton() {
    const selectedTasks = this.checkboxes.filter((box) => box.Checked).length;
    return !this.processed || selectedTasks !== this.complited;
  }

  disableAnalyzeButton() {
    return !this.checkboxes.some((checkbox) => checkbox.Checked) || this.checked !== this.complited;
  }

  childOnChecked(setting: any) {
    const box = this.checkboxes.find((item) => item.BoxId === setting.BoxId);
    if (box) {
      box.Checked = setting.Checked;
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
