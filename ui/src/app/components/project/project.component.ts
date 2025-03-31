import {Component, Input, OnInit} from '@angular/core'
import {IProj} from "../../models/proj.model";
import {ProjectServices} from "../../services/project.services";

@Component({
  selector: 'app-project',
  templateUrl: './project.component.html',
  styleUrls: ['./project.component.css']
})
export class ProjectComponent implements OnInit {
  @Input() project: IProj
  adding: Boolean;

  constructor(private projectService: ProjectServices) {
    //TO_DO
  }

  ngOnInit(): void {
    this.adding = Boolean(this.project.existence);
  }

  addMyProject(project: IProj) {
    if (!this.adding) {
      this.projectService.addProject(String(project.key)).subscribe(resp =>{

      },error => {
        this.adding = !this.adding
        if (error.status == 0){
          alert("Unable to connect to backend")
        }
        if (error.status == 400){
          alert(error.message())
        }
      });
      } else {
      console.log(this.project.id);
        this.projectService.deleteProject(Number(project.id)).subscribe(resp => {

        },
          error => {
            this.adding = !this.adding
            if (error.status == 0){
              alert("Unable to connect to backend")
            }
            if (error.status == 400){
              alert("Unable to connect to DB")
            }
          });
      }
      this.adding = !this.adding
      //TO_DO
    }
}

