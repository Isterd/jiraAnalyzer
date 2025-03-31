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
  error: string | null = null;
  initialized = false; // Заменяем inited на initialized

  constructor(private dbProjectService: DatabaseProjectServices) {}

  ngOnInit(): void {
    this.loadProjects();
    this.initialized = true; // Устанавливаем флаг инициализации
  }

  loadProjects(): void {
    this.loading = true;
    this.error = null;

    this.dbProjectService.getAll().subscribe({
      next: (response) => {
        this.myProjects = response.projects || [];
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading projects:', error);
        this.error = 'Failed to load projects';
        this.loading = false;
      }
    });
  }

  get hasNoProjects(): boolean { // Заменяем noProjects на вычисляемое свойство
    return this.initialized && this.myProjects.length === 0;
  }

  childOnChecked(project: CheckedProject): void {
    if (project.Checked) {
      this.checked.set(project.Name.toString(), project.Id.valueOf());
    } else {
      this.checked.delete(project.Name.toString());
    }
  }
}
