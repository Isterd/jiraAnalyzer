import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

// Импортируйте компоненты для страниц
import { HomePageComponent } from './home-page/home-page.component';
import { ProjectPageComponent } from './project-page/project-page.component';
import { MyProjectPageComponent } from './my-project-page/my-project-page.component';
import { ComparePageComponent } from './compare-page/compare-page.component';
import { CompareProjectPageComponent } from './compare-project-page/compare-project-page.component';

const routes: Routes = [
  { path: '', component: HomePageComponent }, // Главная страница
  { path: 'projects', component: ProjectPageComponent }, // Страница проектов
  { path: 'my-projects', component: MyProjectPageComponent }, // Страница аналитики
  { path: 'compare', component: ComparePageComponent }, // Страница сравнения
  { path: 'compare-projects', component: CompareProjectPageComponent }, // Страница результатов сравнения
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {}
