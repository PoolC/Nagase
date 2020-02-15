import { expect } from 'chai';
import sinon from 'sinon';
import { getRepository } from 'typeorm';
import { Factory } from 'rosie';

import { adminCtx, memberCtx } from '../../fixtures/context';
import Project from '../../../src/models/project';
import ProjectMutation from '../../../src/graphql/mutation/project';
import ProjectService from '../../../src/services/project';

describe('ProjectMutation', () => {
  const service = new ProjectService(getRepository(Project));
  const mutation = new ProjectMutation(service);

  describe('#createProject', () => {
    const project = Factory.build<Project>('Project');
    const dummyInput = {
      name: project.name, body: project.body, genre: project.genre, description: project.description,
      duration: project.duration, participants: project.participants, thumbnailURL: project.thumbnailURL,
    };

    let saveSpy: sinon.SinonStub;
    beforeEach(() => { saveSpy = sinon.stub(ProjectService.prototype, 'save').resolvesArg(0); });
    afterEach(() => saveSpy.restore());

    it('success', async () => {
      const result = await mutation.createProject({ ProjectInput: dummyInput }, adminCtx());
      expect(result.body).to.equal(dummyInput.body);
      expect(saveSpy.calledOnceWith(result)).to.be.true;
    });

    it('failure if permission insufficient', () => {
      expect(() => { mutation.createProject({ ProjectInput: dummyInput }, memberCtx()) }).to.throw(Error, /ERR401/);
    });
  });

  describe('#updateProject', () => {
    const newProject = Factory.build<Project>('Project');
    const dummyInput = {
      name: newProject.name, body: newProject.body, genre: newProject.genre, description: newProject.description,
      duration: newProject.duration, participants: newProject.participants, thumbnailURL: newProject.thumbnailURL,
    };

    let findProjectSpy: sinon.SinonStub;
    let saveProjectSpy: sinon.SinonStub;
    let project: Project;
    beforeEach(() => {
      findProjectSpy = sinon.stub(ProjectService.prototype, 'findById').resolves(project);
      saveProjectSpy = sinon.stub(ProjectService.prototype, 'save').resolvesArg(0);
    });
    afterEach(() => [findProjectSpy, saveProjectSpy].forEach((each) => each.restore()));

    context('project are valid', () => {
      before(() => { project = Factory.build<Project>('Project'); });

      it('success if permission sufficient', async () => {
        const result = await mutation.updateProject({ projectID: project.id, ProjectInput: dummyInput }, adminCtx());
        expect(result.body).to.equal(dummyInput.body);
        expect(saveProjectSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail if permission insufficient', () => {
        expect(() => { mutation.createProject({ ProjectInput: dummyInput }, memberCtx()) }).to.throw(Error, /ERR401/);
        expect(saveProjectSpy.notCalled).to.be.true;
      });
    });

    context('project not exists', () => {
      before(() => { project = null });
      it('fail', () => {
        expect(mutation.updateProject({ projectID: 0, ProjectInput: dummyInput }, adminCtx())).to.be.rejectedWith(Error, /ERR400/);
        expect(saveProjectSpy.notCalled).to.be.true;
      });
    });
  });

  describe('#deleteProject', () => {
    let findProjectSpy: sinon.SinonStub;
    let deleteProjectSpy: sinon.SinonStub;
    let project: Project;
    beforeEach(() => {
      findProjectSpy = sinon.stub(ProjectService.prototype, 'findById').resolves(project);
      deleteProjectSpy = sinon.stub(ProjectService.prototype, 'delete').resolvesArg(0);
    });
    afterEach(() => [findProjectSpy, deleteProjectSpy].forEach((each) => each.restore()));

    context('member and project are valid', () => {
      before(() => { project = Factory.build<Project>('Project'); });

      it('success if permission sufficient', async () => {
        const result = await mutation.deleteProject({ projectID: project.id }, adminCtx());
        expect(result).not.to.be.null;
        expect(deleteProjectSpy.calledOnceWith(result)).to.be.true;
      });

      it('fail if permission insufficient', () => {
        expect(() => { mutation.deleteProject({ projectID: project.id }, memberCtx()) }).to.throw(Error, /ERR401/);
        expect(deleteProjectSpy.notCalled).to.be.true;
      });
    });

    context('project not exists', () => {
      before(() => { project = null });
      it('fail', async () => {
        expect(await mutation.deleteProject({ projectID: 0 }, adminCtx())).to.be.null;
        expect(deleteProjectSpy.notCalled).to.be.true;
      });
    });
  });
});
