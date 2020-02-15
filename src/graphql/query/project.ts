import { Service } from 'typedi';
import { Context } from 'koa';

import Project from '../../models/project';
import ProjectService from '../../services/project';

export interface ProjectIDInput {
  projectID: number;
}

@Service()
export default class ProjectQuery {
  constructor(
    private readonly projectService: ProjectService,
  ) {}

  public async project(args: ProjectIDInput, _: Partial<Context>): Promise<Project> {
    return this.projectService.findById(args.projectID);
  }

  public async projects(_1: object, _2: Partial<Context>): Promise<Project[]> {
    return (await this.projectService.findAll()).sort((a, b) => b.id - a.id);
  }
}
