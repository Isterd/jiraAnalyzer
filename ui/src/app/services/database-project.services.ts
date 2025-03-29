import { Injectable } from "@angular/core";
import { HttpClient } from "@angular/common/http";
import {catchError, Observable, tap} from "rxjs";
import { IRequest } from "../models/request.model"; // Используем только IRequest
import { ConfigurationService } from "./configuration.services";

@Injectable({
  providedIn: 'root'
})
export class DatabaseProjectServices {
  urlPath = "";

  constructor(private http: HttpClient, private configurationService: ConfigurationService) {
    this.urlPath = configurationService.getValue("pathUrl");
  }

  // Получение всех проектов с пагинацией и фильтрацией
  getAll(): Observable<IRequest> {
    const url = `${this.urlPath}/api/v1/projects`;
    return this.http.get<IRequest>(url);
  }

  // Получение аналитики для конкретного проекта
  getProjectAnalytics(projectKey: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/projects/${projectKey}/analytics`;
    return this.http.get<any>(url).pipe(
      tap((response) => {
        console.log("Ответ сервера:", response);
      }),
      catchError((error) => {
        console.error("Ошибка при получении аналитики:", error);
        throw error;
      })
    );
  }

  // Получение графика по задаче (GET)
  getGraph(taskNumber: string, projectName: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/graph/get/${taskNumber}?project=${projectName}`;
    return this.http.get<any>(url);
  }

  // Создание графика по задаче (POST)
  makeGraph(taskNumber: string, projectName: string): Observable<any> {
    const url = `${this.urlPath}/api/v1/graph/make/${taskNumber}?project=${projectName}`;
    return this.http.post<any>(url, {});
  }

  // Получение сравнения двух проектов (GET)
  getComparison(taskNumber: string, projectKeys: string[]): Observable<any> {
    const url = `${this.urlPath}/api/v1/compare/${taskNumber}?project=${projectKeys.join(',')}`;
    return this.http.get<any>(url);
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
