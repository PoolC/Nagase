import { expect } from 'chai';
import sinon from 'sinon';
import { Factory } from 'rosie';
import { getRepository } from 'typeorm';

import ProjectQuery from '../../../src/graphql/query/project';
import Project from '../../../src/models/project';
import ProjectService from '../../../src/services/project';

describe('ProjectQuery', () => {
  const query = new ProjectQuery(new ProjectService(getRepository(Project)));

  describe('#project', () => {
    let project: Project;
    let findProjectSpy: sinon.SinonStub;
    beforeEach(() => { findProjectSpy = sinon.stub(ProjectService.prototype, 'findById').resolves(project) });
    afterEach(() => findProjectSpy.restore());

    context('project exists', () => {
      before(() => { project = Factory.build<Project>('Project'); });
      it('success', async () => {
        const result = await query.project({ projectID: project.id }, null);
        expect(result.name).to.equal(project.name);
        expect(result.body).to.equal(project.body);
      });
    });

    context('project not exists', () => {
      before(() => { project = null; });
      it('return null', async () => expect(await query.project({ projectID: -1 }, null)).to.be.null);
    });
  });

  describe('#projects', () => {
    let findProjectsSpy: sinon.SinonStub;
    beforeEach(() => { findProjectsSpy = sinon.stub(ProjectService.prototype, 'findAll').resolves(Factory.buildList('Project', 10)); });
    afterEach(() => findProjectsSpy.restore());

    it('success', async () => {
      const results = await query.projects({}, null);
      expect(results.length).to.equal(10);
      for (let i = 0; i < 9; i++) {
        expect(results[i].id).to.greaterThan(results[i + 1].id);
      }
    });
  });
});
