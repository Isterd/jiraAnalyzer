import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { IProj } from "../../models/proj.model";
import { CheckedProject } from "../../models/check-element.model";

@Component({
  selector: 'app-checkbox-with-project', // Изменяем селектор для единообразия
  templateUrl: './checkbox-with-project.component.html',
  styleUrls: ['./checkbox-with-project.component.css']
})
export class ProjectWithCheckboxComponent implements OnInit {
  @Output() onChecked = new EventEmitter<CheckedProject>();
  @Input() project!: IProj;
  isChecked: boolean;

  constructor() {}

  ngOnInit(): void {
    this.isChecked = this.project.existence ?? false;
  }

  changed(event: Event): void {
    const isChecked = (event.target as HTMLInputElement).checked;
    this.onChecked.emit({
      Name: this.project.name,
      Checked: isChecked,
      Id: this.project.id
    });
  }
}
