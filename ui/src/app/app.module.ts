import { APP_INITIALIZER, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { NgxPaginationModule } from 'ngx-pagination';
import { ChartModule } from 'angular-highcharts';

// Компоненты
import { AppComponent } from './app.component';
import { HomePageComponent } from './home-page/home-page.component';
import { ProjectComponent } from './components/project/project.component';
import { ProjectPageComponent } from './project-page/project-page.component';
import { MyProjectPageComponent } from './my-project-page/my-project-page.component';
import { MyProjectComponent } from './components/my-project/my-project.component';
import { ComparePageComponent } from './compare-page/compare-page.component';
import { ProjectWithCheckboxComponent } from './components/checkbox-with-project/checkbox-with-project.component';
import { CompareProjectPageComponent } from './compare-project-page/compare-project-page.component';
import { CheckboxWithSettingsComponent } from './components/checkbox-with-settings/checkbox-with-settings.component';
import { ProjectStatPageComponent } from './project-stat-page/project-stat-page.component';

// Сервисы
import { ConfigurationService } from './services/configuration.services';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

// Функция инициализации приложения
export function initApp(configurationService: ConfigurationService) {
  return () => configurationService.load().toPromise();
}

// Маршруты
const routes = [
  { path: '', component: HomePageComponent }, // Главная страница
  { path: 'projects', component: ProjectPageComponent }, // Страница проектов
  { path: 'compare', component: ComparePageComponent }, // Страница сравнения
  { path: 'myprojects', component: MyProjectPageComponent }, // Страница аналитики
  { path: 'compare-projects', component: CompareProjectPageComponent }, // Результаты сравнения
  { path: 'projects-settings', component: CheckboxWithSettingsComponent }, // Настройки проектов
  { path: 'project-stat/:id', component: ProjectStatPageComponent } // Статистика проекта (динамический параметр)
];

@NgModule({
  declarations: [
    AppComponent,
    HomePageComponent,
    ProjectComponent,
    ProjectPageComponent,
    MyProjectComponent,
    MyProjectPageComponent,
    ComparePageComponent,
    ProjectWithCheckboxComponent,
    CompareProjectPageComponent,
    CheckboxWithSettingsComponent,
    ProjectStatPageComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule, // Для HTTP-запросов
    FormsModule, // Для двухсторонней привязки данных
    NgxPaginationModule, // Для пагинации
    ChartModule, // Для графиков Highcharts
    RouterModule.forRoot(routes), BrowserAnimationsModule // Маршрутизация
  ],
  providers: [
    {
      provide: APP_INITIALIZER,
      useFactory: initApp,
      multi: true,
      deps: [ConfigurationService]
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {}
