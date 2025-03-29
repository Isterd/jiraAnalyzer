import { Component, OnInit } from '@angular/core';
import { IProj } from "../models/proj.model";
import { DatabaseProjectServices } from "../services/database-project.services";
import { Router } from '@angular/router';
import { CheckedProject } from "../models/check-element.model";
import { ConfigurationService } from "../services/configuration.services";

@Component({
  selector: 'app-compare-page',
  templateUrl: './compare-page.component.html',
  styleUrls: ['./compare-page.component.css']
})
export class ComparePageComponent implements OnInit {
  projects: IProj[] = [];
  checked: Map<string, number> = new Map(); // Уточнение типов
  noProjects: boolean = false;
  inited: boolean = false;

  webUrl: string = ""; // Добавлен явный тип

  constructor(
    private configurationService: ConfigurationService,
    private myProjectService: DatabaseProjectServices,
    private router: Router
  ) {
    this.webUrl = configurationService.getValue("webUrl");
  }

  ngOnInit(): void {
    this.myProjectService.getAll().subscribe(
      (projects) => {
        if (projects && projects.projects) {
          this.noProjects = projects.projects.length === 0;
          this.projects = projects.projects;
          this.inited = true;
        } else {
          console.error('Ответ сервера некорректен:', projects);
          this.noProjects = true;
        }
      },
      (error) => {
        console.error('Ошибка при загрузке проектов:', error);
        alert('Не удалось загрузить проекты. Проверьте соединение с сервером.');
      }
    );
  }

  onClickCompare(): void {
    const items: string[] = [];
    const ids: number[] = [];

    this.checked.forEach((value: number, key: string) => {
      if (value) {
        items.push(key);
        ids.push(value);
      }
    });

    if (items.length > 3) {
      this.showErrorMessage("Максимальное число проектов для сравнения — 3.");
    } else if (items.length < 2) {
      this.showErrorMessage("Минимальное число проектов для сравнения — 2.");
    } else {
      this.router.navigate(['/compare-projects'], {
        queryParams: {
          keys: items.join(','), // Преобразуем массив в строку
          values: ids.join(',') // Преобразуем массив в строку
        }
      });
    }
  }

  childOnChecked(project: CheckedProject): void {
    if (project.Checked) {
      this.checked.set(project.Name.toString(), project.Id.valueOf()); // Преобразование Number -> number
    } else {
      this.checked.delete(project.Name.toString());
    }
  }

  showErrorMessage(msg: string): void {
    alert(msg);
  }
}
