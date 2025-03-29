import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {catchError, Observable, throwError} from 'rxjs';
import { IRequest } from '../models/request.model';
import { ConfigurationService } from './configuration.services';

@Injectable({
  providedIn: 'root'
})
export class ProjectServices {
  private urlPath = '';

  constructor(private http: HttpClient, private configurationService: ConfigurationService) {
    this.urlPath = configurationService.getValue("pathUrl");
  }

  // Получение всех проектов с учетом пагинации и поиска
  getAll(page: number, searchName: string): Observable<IRequest> {
    const url = `${this.urlPath}/api/v1/projects?page=${page}&limit=10&search=${searchName}`;
    console.log('Запрос к URL:', url); // Логирование URL
    return this.http.get<IRequest>(url).pipe(
      catchError((error) => {
        console.error('Ошибка при загрузке проектов:', error);
        return throwError(() => new Error('Не удалось загрузить проекты. Попробуйте позже.'));
      })
    );
  }

  // Добавление проекта в базу данных
  addProject(key: string): Observable<IRequest> {
    const url = `${this.urlPath}/api/v1/connector/updateProject`;
    return this.http.post<IRequest>(url, { project: key }).pipe(
      catchError((error) => {
        console.error('Ошибка при добавлении проекта:', error);
        return throwError(() => new Error('Не удалось добавить проект. Попробуйте позже.'));
      })
    );
  }

  // Удаление проекта из базы данных
  deleteProject(id: number): Observable<IRequest> {
    const url = `${this.urlPath}/api/v1/projects/${id}`;
    return this.http.delete<IRequest>(url).pipe(
      catchError((error) => {
        console.error('Ошибка при удалении проекта:', error);
        return throwError(() => new Error('Не удалось удалить проект. Попробуйте позже.'));
      })
    );
  }
}
