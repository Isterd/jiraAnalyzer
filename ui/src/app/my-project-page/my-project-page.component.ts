import { Component, OnInit } from '@angular/core';
import { DatabaseProjectServices } from "../services/database-project.services";
import { IProj } from "../models/proj.model";
import { CheckedProject } from "../models/check-element.model";

@Component({
  selector: 'app-my-project-page',
  templateUrl: './my-project-page.component.html',
  styleUrls: ['./my-project-page.component.css']
})
export class MyProjectPageComponent implements OnInit {
  myProjects: IProj[] = [];
  checked: Map<string, number> = new Map(); // Уточнение типов
  loading = false;
  noProjects = false;
  inited = false;

  constructor(private myProjectService: DatabaseProjectServices) {}

  ngOnInit(): void {
    this.loading = true;
    this.myProjectService.getAll().subscribe(
      (projects) => {
        if (projects && projects.projects) {
          this.noProjects = projects.projects.length === 0;
          this.myProjects = projects.projects;
          this.inited = true;
        } else {
          console.error('Ответ сервера некорректен:', projects);
          this.noProjects = true;
        }
        this.loading = false;
      },
      (error) => {
        console.error('Ошибка при загрузке проектов:', error);
        alert('Не удалось загрузить проекты. Проверьте соединение с сервером.');
        this.loading = false;
      }
    );
  }

  childOnChecked(project: CheckedProject): void {
    if (project.Checked) {
      this.checked.set(project.Name.toString(), project.Id.valueOf());
    } else {
      this.checked.delete(project.Name.toString());
    }
  }
}
