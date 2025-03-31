export interface AnalyticsApiResponse {
  data: {
    total_issues: number;
    closed_issues: number;
    open_issues: number;
    reopen_issues: number;
    resolved_issues: number;
    in_progress_issues: number;
    average_time_issues: number;
    average_count_issues: number;
  };
}

export interface AnalyticsFrontendData {
  data: {
    allIssuesCount: number;
    openIssuesCount: number;
    closeIssuesCount: number;
    reopenedIssuesCount: number;
    resolvedIssuesCount: number;
    progressIssuesCount: number;
    averageTime: number;
    averageIssuesCount: number;
  };
}
