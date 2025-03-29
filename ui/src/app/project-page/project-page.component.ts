import { Component, OnInit } from '@angular/core';
import { ProjectServices } from '../services/project.services';
import { IProj } from '../models/proj.model';
import { PageInfo } from '../models/pageInfo.model';

@Component({
  selector: 'app-project-page',
  templateUrl: './project-page.component.html',
  styleUrls: ['./project-page.component.css']
})
export class ProjectPageComponent implements OnInit {
  projects: IProj[] = [];
  loading = false;
  searchName = '';
  pageInfo: PageInfo = { currentPage: 1, pageCount: 0, totalCount: 0 };
  start_page = 1;

  constructor(private projectService: ProjectServices) {}

  ngOnInit(): void {
    this.loadProjects();
  }

  loadProjects(): void {
    this.loading = true;
    this.projectService.getAll(this.start_page, this.searchName).subscribe(
      (response) => {
        if (response && response.projects) {
          this.projects = response.projects;
          this.pageInfo = response.pageInfo;
        } else {
          console.error('Ответ сервера пуст или некорректен:', response);
          this.projects = [];
          this.pageInfo = { currentPage: 1, pageCount: 0, totalCount: 0 };
        }
        this.loading = false;
      },
      (error) => {
        this.loading = false;
        console.error('Ошибка при загрузке проектов:', error);
        alert('Не удалось загрузить проекты. Попробуйте позже.');
      }
    );
  }

  gty(page: number): void {
    this.loading = true;
    this.projectService.getAll(page, this.searchName).subscribe(
      (response) => {
        this.projects = response.projects;
        this.pageInfo = response.pageInfo;
        this.loading = false;
      },
      (error) => {
        this.loading = false;
        console.error('Ошибка при загрузке проектов:', error);
        alert('Не удалось загрузить проекты. Попробуйте позже.');
      }
    );
  }

  // Метод для поиска проектов
  getSearchProjects(): void {
    this.pageInfo.currentPage = this.start_page; // Сбрасываем страницу на первую
    this.gty(this.pageInfo.currentPage); // Загружаем проекты с новым поисковым запросом
  }

  // Обработчик события пагинации
  onPageChange(page: number): void {
    this.pageInfo.currentPage = page;
    this.gty(page);
  }
}
