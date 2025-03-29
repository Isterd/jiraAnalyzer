import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { DatabaseProjectServices } from "../services/database-project.services";
import { Chart } from "angular-highcharts";
import { openTaskChartOptions } from "./helpers/openTaskChartOptions";
import { ConfigurationService } from "../services/configuration.services";

@Component({
  selector: 'app-compare-project-page',
  templateUrl: './compare-project-page.component.html',
  styleUrls: ['./compare-project-page.component.scss']
})
export class CompareProjectPageComponent implements OnInit {
  projects: string[] = [];
  ids: string[] = [];
  resultReq: any[] = []; // Используем any для данных проекта
  openTaskChart = new Chart();
  webUrl: string = ""; // Добавляем свойство webUrl

  constructor(
    private route: ActivatedRoute,
    private dbProjectService: DatabaseProjectServices,
    private configurationService: ConfigurationService // Добавляем сервис конфигурации
  ) {
    this.projects = this.route.snapshot.queryParamMap.getAll("keys");
    this.ids = this.route.snapshot.queryParamMap.getAll("values");
    this.webUrl = configurationService.getValue("webUrl"); // Инициализируем webUrl
  }

  ngOnInit(): void {
    for (let i = 0; i < this.projects.length; i++) {
      this.dbProjectService.getProjectAnalytics(this.ids[i]).subscribe(
        (projects: any) => { // Указываем явный тип для projects
          this.resultReq[i] = projects.data;
        },
        (error: any) => { // Указываем явный тип для error
          console.error(`Ошибка при получении статистики для проекта ${this.projects[i]}:`, error);
        }
      );
    }

    const colors = ["blue", "green", "red", "orange", "purple", "black"];

    this.dbProjectService.getComparison("5", this.projects).subscribe(
      (info: any) => { // Указываем явный тип для info
        if (info.data["count"] == null) {
          const openTaskElem = document.getElementById('open-task');
          const openTaskTitle = document.getElementById('open-task-title');
          if (openTaskElem) openTaskElem.remove();
          if (openTaskTitle) openTaskTitle.remove();
        } else {
          if (openTaskChartOptions.xAxis && 'categories' in openTaskChartOptions.xAxis) {
            openTaskChartOptions.xAxis.categories = info.data["categories"];
          }

          for (let j = 0; j < this.projects.length; j++) {
            const count = [];
            for (let i = 0; i < info.data["categories"].length; i++) {
              count.push(info.data["count"][info.data["categories"][i]][j]);
            }
            openTaskChartOptions.series?.push({
              name: this.projects[j],
              type: "column",
              color: colors[j],
              data: count
            });
          }
          this.openTaskChart = new Chart(openTaskChartOptions);
        }
      },
      (error: any) => { // Указываем явный тип для error
        console.error('Ошибка при получении графика задач:', error);
      }
    );
  }

  ngOnDestroy(): void {
    if (openTaskChartOptions.xAxis && 'categories' in openTaskChartOptions.xAxis) {
      openTaskChartOptions.xAxis.categories = [];
    }
    openTaskChartOptions.series = [];
  }
}
