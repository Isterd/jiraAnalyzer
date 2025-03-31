export interface IComparisonResponse {
  _links: {
    self: string;
  };
  data: IProjectComparison[];
}

export interface IProjectComparison {
  project_key: string;
  data: IMetric[];
}

export interface IMetric {
  Priority?: string; // Для taskNumber=5
  TaskCount: number;
  Day?: string; // Для taskNumber=3
  CumulativeOpened?: number;
  CumulativeClosed?: number;
}
