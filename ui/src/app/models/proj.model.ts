export interface IProj {
  id: number;
  key: string;
  name: string;
  url: string; // Добавляем отсутствующее свойство
  existence?: boolean;
  total_issues?: number;
  open_issues?: number;
  closed_issues?: number;
  reopened_issues?: number;
  resolved_issues?: number;
  in_progress_issues?: number;
  average_time?: number;
  average_issues_per_day?: number;
}
