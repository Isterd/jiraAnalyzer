import { Component, OnInit } from '@angular/core';
import { DatabaseProjectServices } from '../services/database-project.services';
import { Router } from '@angular/router';
import { IProj } from '../models/proj.model';

@Component({
  selector: 'app-compare-page',
  templateUrl: './compare-page.component.html',
  styleUrls: ['./compare-page.component.css']
})
export class ComparePageComponent implements OnInit {
  projects: IProj[] = [];
  selectedProjectKeys = new Set<string>(); // Используем Set для хранения ключей
  isLoading = true;
  errorMessage: string | null = null;

  constructor(
    private projectService: DatabaseProjectServices,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.loadProjects();
  }

  loadProjects(): void {
    this.projectService.getAll().subscribe({
      next: (response: any) => {
        console.log('Получены проекты:', response);
        this.projects = response.projects || [];
        this.isLoading = false;

        if (this.projects.length === 0) {
          this.errorMessage = 'Нет доступных проектов для сравнения';
        }
      },
      error: (err) => {
        console.error('Ошибка загрузки проектов:', err);
        this.errorMessage = 'Ошибка загрузки проектов';
        this.isLoading = false;
      }
    });
  }

  toggleProjectSelection(project: IProj, event: Event): void {
    const checkbox = event.target as HTMLInputElement;
    const isChecked = checkbox.checked;

    if (isChecked) {
      if (this.selectedProjectKeys.size >= 3) {
        checkbox.checked = false;
        this.errorMessage = 'Можно выбрать не более 3 проектов';
        return;
      }
      this.selectedProjectKeys.add(project.key);
    } else {
      this.selectedProjectKeys.delete(project.key);
    }

    this.errorMessage = null;
    console.log('Выбранные проекты:', Array.from(this.selectedProjectKeys));
  }

  isSelected(project: IProj): boolean {
    return this.selectedProjectKeys.has(project.key);
  }

  compareProjects(): void {
    const selectedCount = this.selectedProjectKeys.size;
    console.log('Попытка сравнения. Выбрано проектов:', selectedCount);

    if (selectedCount < 2 || selectedCount > 3) {
      this.errorMessage = 'Выберите от 2 до 3 проектов для сравнения';
      return;
    }

    const keys = Array.from(this.selectedProjectKeys).join(',');
    const selectedProjects = this.projects.filter(p => this.selectedProjectKeys.has(p.key));
    const ids = selectedProjects.map(p => p.key).join(','); // Используем key как идентификатор

    console.log('Параметры для сравнения:', { keys, ids });

    this.router.navigate(['/compare-projects'], {
      queryParams: { keys, values: ids }
    }).catch(err => {
      console.error('Ошибка навигации:', err);
      this.errorMessage = 'Ошибка перехода на страницу сравнения';
    });
  }
}
