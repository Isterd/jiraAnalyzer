import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import {catchError, map, Observable, tap, throwError} from "rxjs";
import { IRequest } from "../models/request.model"; // Используем только IRequest
import { ConfigurationService } from "./configuration.services";
import {IComparisonResponse} from "../models/comparison.model";


@Injectable({
  providedIn: 'root'
})
export class DatabaseProjectServices {
  urlPath = "";

  constructor(private http: HttpClient, private configurationService: ConfigurationService) {
    this.urlPath = configurationService.getValue("pathUrl");
  }

  // Получение всех проектов с пагинацией и фильтрацией
  getAll(): Observable<any> {
    const url = `${this.urlPath}/api/v1/projects`;
    console.log('Отправка запроса на:', url); // Добавим лог URL

    return this.http.get(url).pipe(
      tap(response => console.log('Сырой ответ сервера:', response)),
      catchError(error => {
        console.error('Ошибка при загрузке проектов:', error);
        return throwError(() => error);
      })
    );
  }

  // Добавим метод для получения статистики проекта
  getProjectStats(projectId: number): Observable<any> {
    return this.getProjectAnalytics(projectId.toString()).pipe(
      map(response => ({
        allIssuesCount: response.data.allIssuesCount,
        openIssuesCount: response.data.openIssuesCount,
        closeIssuesCount: response.data.closeIssuesCount,
        reopenedIssuesCount: response.data.reopenedIssuesCount,
        resolvedIssuesCount: response.data.resolvedIssuesCount,
        progressIssuesCount: response.data.progressIssuesCount,
        averageTime: response.data.averageTime,
        averageIssuesCount: response.data.averageIssuesCount
      })),
      catchError(error => {
        console.error('Error loading project stats:', error);
        return throwError(() => new Error('Не удалось загрузить статистику проекта'));
      })
    );
  }

  // Получение аналитики для конкретного проекта
  getProjectAnalytics(projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/projects/${projectKey}/analytics`;
    console.log('Запрос аналитики проекта:', url);

    return this.http.get(url).pipe(
      map((response: any) => {
        console.log('Сырой ответ аналитики:', response);

        // Правильное преобразование данных с учетом имен полей из API
        return {
          data: {
            allIssuesCount: response.data?.total_issues ?? 0,
            openIssuesCount: response.data?.open_issues ?? 0,
            closeIssuesCount: response.data?.closed_issues ?? 0,
            reopenedIssuesCount: response.data?.reopen_issues ?? 0, // было reopened_issues
            resolvedIssuesCount: response.data?.resolved_issues ?? 0,
            progressIssuesCount: response.data?.in_progress_issues ?? 0,
            averageTime: response.data?.average_time_issues ?? 0, // было average_time
            averageIssuesCount: response.data?.average_count_issues ?? 0 // было average_issues_per_day
          }
        };
      }),
      catchError(error => {
        console.error('Ошибка запроса аналитики:', error);
        return throwError(() => error);
      })
    );
  }

  // Получение графика по задаче (GET)
  getGraph(taskNumber: string, projectName: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/graph/get/${taskNumber}?project=${projectName}`;
    return this.http.get<any>(url);
  }

  // Создание графика по задаче (POST)
  makeGraph(taskNumber: string, projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/graph/make/${taskNumber}?project=${projectKey}`;
    console.log('Отправка запроса на создание графика:', url);

    return this.http.post(url, {}).pipe(
      map((response: any) => {
        if (!response.success) {
          throw new Error(response.message || 'Неизвестная ошибка сервера');
        }
        return response;
      }),
      catchError((error) => {
        console.error('Ошибка при создании графика:', error);
        return throwError(() => new Error(error.error?.message || 'Ошибка сервера'));
      })
    );
  }

  deleteGraphs(projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/graph/delete/${projectKey}`;
    return this.http.delete<any>(url, {})
  }

  isAnalyzed(projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/isAnalyzed/${projectKey}`;
    return this.http.get<any>(url);
  }

  // Получение сравнения двух проектов (GET)
  getComparison(taskNumber: string, projectKeys: string[]): Observable<IComparisonResponse> {
    const url = `${this.urlPath}/api/v1/compare/${taskNumber}?project=${projectKeys.join(',')}`;
    return this.http.get<IComparisonResponse>(url).pipe(
      catchError((error) => {
        console.error('Error fetching comparison:', error);
        throw error;
      })
    );
  }

  // Создание сравнения двух проектов (POST)
  makeComparison(taskNumber: string, projectKeys: string[]): Observable<any> {
    const url = `${this.urlPath}/api/v1/compare/${taskNumber}?project=${projectKeys.join(',')}`;
    return this.http.post<any>(url, {});
  }

  // Обновление проекта через коннектор
  updateConnectorProject(projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/connector/updateProject`;
    return this.http.post<any>(url, { project: projectKey });
  }

  // Получение проектов через коннектор
  getConnectorProjects(page: number, limit: number, search: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/connector/projects?limit=${limit}&page=${page}&search=${search}`;
    return this.http.get<any>(url);
  }
}
