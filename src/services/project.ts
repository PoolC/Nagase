import { Service } from 'typedi';
import { Repository } from 'typeorm';
import { InjectRepository } from 'typeorm-typedi-extensions';

import Project from '../models/project';

@Service()
export default class ProjectService {
  constructor(
    @InjectRepository(Project) private readonly projectRepository: Repository<Project>,
  ) {}

  public async findById(projectId: number): Promise<Project> {
    return this.projectRepository.findOne({ id: projectId });
  }

  public async findAll(): Promise<Project[]> {
    return this.projectRepository.find();
  }

  public async save(obj: Project): Promise<Project> {
    return this.projectRepository.save(obj);
  }

  public async delete(obj: Project): Promise<void> {
    await this.projectRepository.delete(obj);
  }
}
