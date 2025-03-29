import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {mapTo, Observable, tap} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ConfigurationService {
  private configuration: any;

  constructor(private http: HttpClient) {}

  load(): Observable<void> {
    return this.http.get('/assets/config.json').pipe(
      tap((config: any) => this.configuration = config),
      mapTo(undefined)
    );
  }

  getValue(key: string): string {
    return this.configuration[key];
  }
}
