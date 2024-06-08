<template>
  <div class="container mt-5">
    <h2 class="text-center">My Repositories</h2>
    <div class="row">
      <div v-for="repo in repositories" :key="repo.id" class="col-md-4 mb-3">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ repo.name }} <span class="badge" :class="repo.is_public ? 'text-bg-success' : 'text-bg-secondary'">{{ repo.is_public ? 'public' : 'private' }}</span></h5>
            <p class="card-text">URL: http://localhost:8000/{{ repo.username }}/{{ repo.name }}</p>
            <div class="dropdown">
              <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
                Actions
              </button>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
                <li><a class="dropdown-item" href="#" @click="viewCommits(repo.id)">View Commits</a></li>
                <li><a class="dropdown-item" href="#" @click="deleteRepo(repo.id)">Delete Repo</a></li>
                <li><a class="dropdown-item" href="#" @click="viewContents(repo.id)">View Contents</a></li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal для отображения содержимого репозитория -->
    <div class="modal fade" id="repoContentsModal" tabindex="-1" aria-labelledby="repoContentsModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="repoContentsModalLabel">Repository Contents</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <ul>
              <li v-for="file in repoContents" :key="file">{{ file }}</li>
            </ul>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'MyRepositories',
  data() {
    return {
      repositories: [],
      repoContents: [], // Для хранения содержимого репозитория
    };
  },
  methods: {
    async fetchRepositories() {
      try {
        const response = await axios.post('http://localhost:8000/api/v1/user_repos', { id: 1 }, { withCredentials: true });
        this.repositories = response.data;
      } catch (error) {
        console.error('Error fetching repositories:', error);
      }
    },
    async deleteRepo(repoId) {
      try {
        await axios.post('http://localhost:8000/api/v1/delete_repo', { repo_id: repoId }, { withCredentials: true });
        this.fetchRepositories(); // Обновить список репозиториев после удаления
      } catch (error) {
        console.error('Error deleting repository:', error);
      }
    },
    async viewContents(repoId) {
      try {
        const response = await axios.post('http://localhost:8000/api/v1/view_contents', { repo_id: repoId }, { withCredentials: true });
        this.repoContents = response.data.contents;
        new bootstrap.Modal(document.getElementById('repoContentsModal')).show();
      } catch (error) {
        console.error('Error fetching repository contents:', error);
      }
    }
  },
  created() {
    this.fetchRepositories();
  }
};
</script>

<style scoped>
.card {
  margin-bottom: 20px;
}
</style>
