import { Context } from 'koa';
import { Service } from 'typedi';

import Project from '../../models/project';
import ProjectService from '../../services/project';
import { Permission, PermissionGuard } from '../guard';

export interface CreateProjectInput {
  ProjectInput: {
    name: string;
    body: string;
    genre: string;
    description: string;
    duration: string;
    participants: string;
    thumbnailURL: string;
  };
}

export interface UpdateProjectInput {
  projectID: number;
  ProjectInput: {
    name: string;
    body: string;
    genre: string;
    description: string;
    duration: string;
    participants: string;
    thumbnailURL: string;
  };
}

export interface DeleteProjectInput {
  projectID: number;
}

@Service()
export default class ProjectMutation {
  constructor(
    private readonly projectService: ProjectService,
  ) {}

  @PermissionGuard(Permission.Admin)
  public async createProject(args: CreateProjectInput, ctx: Partial<Context>): Promise<Project> {
    const project = new Project();
    project.name = args.ProjectInput.name;
    project.body = args.ProjectInput.body;
    project.genre = args.ProjectInput.genre;
    project.description = args.ProjectInput.description;
    project.duration = args.ProjectInput.duration;
    project.participants = args.ProjectInput.participants;
    project.thumbnailURL = args.ProjectInput.thumbnailURL;

    return this.projectService.save(project);
  }

  @PermissionGuard(Permission.Admin)
  public async updateProject(args: UpdateProjectInput, ctx: Partial<Context>): Promise<Project> {
    const project = await this.projectService.findById(args.projectID);
    if (!project) {
      throw new Error('ERR400');
    }

    project.name = args.ProjectInput.name || project.name;
    project.body = args.ProjectInput.body || project.body;
    project.genre = args.ProjectInput.genre || project.genre;
    project.description = args.ProjectInput.description || project.description;
    project.duration = args.ProjectInput.duration || project.duration;
    project.participants = args.ProjectInput.participants || project.participants;
    project.thumbnailURL = args.ProjectInput.thumbnailURL || project.thumbnailURL;

    return this.projectService.save(project);
  }

  @PermissionGuard(Permission.Admin)
  public async deleteProject(args: DeleteProjectInput, ctx: Partial<Context>): Promise<Project> {
    const project = await this.projectService.findById(args.projectID);
    if (project) {
      await this.projectService.delete(project);
    }

    return project;
  }
}
