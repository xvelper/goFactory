<template>
  <div class="container mt-5">
    <div class="d-flex justify-content-between">
      <div> </div>
      <div class="repoName">Мои репозитории</div>
      <button type="button" class="btn btn-primary" data-bs-toggle="modal" data-bs-target="#repoAddModal">Добавить <br> репозиторий</button>
    </div>
    <div class="row">
      <div v-for="repo in repositories" :key="repo.ID" class="col-md-4 mb-3">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ repo.Name }} <span class="badge text-bg-secondary">{{ repo.IsPublic ? 'public' : 'private' }}</span></h5>
            <p class="card-text">URL: https://api.xvelper.ru/{{ repo.username }}/{{ repo.Name }}</p>
            <div class="dropdown">
              <button class="btn btn-secondary dropdown-toggle" type="button" id="dropdownMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
                Действие
              </button>
              <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton">
                <li><a class="dropdown-item" href="#" @click="viewRepoContent(repo.ID, repo.Name)">Посмотреть содержимое</a></li>
                <li><a class="dropdown-item" data-bs-toggle="modal" href="#repoViewCommits" @click="viewCommits(repo.ID)">Посмотреть коммиты</a></li>
                <li><a class="dropdown-item" href="#" @click="deleteRepo(repo.ID, repo.OwnerID)">Удалить репозиторий</a></li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="modal" id="repoAddModal" tabindex="-1">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Добавление репозитория</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div class="mb-3">
              <label for="repoName" class="form-label">Название репозитория</label>
              <input type="text" class="form-control" id="repoName" v-model="newRepoName" placeholder="Пример: test, gogs">
            </div>
            <div class="mb-3 form-check">
              <input type="checkbox" v-model="isPublic" class="form-check-input" id="isPublic">
              <label class="form-check-label" for="isPublic">Публичный репозиторий</label>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
            <button type="button" class="btn btn-primary" @click="createRepo" data-bs-dismiss="modal">Добавить</button>
          </div>
        </div>
      </div>
    </div>

    <div class="modal fade" id="repoViewCommits" tabindex="-1" aria-labelledby="repoViewCommitsLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="repoViewCommitsLabel">Просмотр коммитов репозитория</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <div v-for="commit in commitLogs" :key="commit" class="commit-item mb-3">
              <div class="commits">
                <div class="commits-head">
                  <h6 class="mb-1">{{ commit.Message }}</h6>
                  <small class="text-muted">{{ commit.Hash }}</small>
                </div>
                <div class="commits-body">
                  <small class="text-muted">{{ commit.Author }} committed {{ commit.Date }}</small>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Закрыть</button>
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
      newRepoName: '',
      isPublic: false,
      commitLogs: [],
    };
  },
  methods: {
    async fetchRepositories() {
      try {
        const response = await axios.post('https://api.xvelper.ru/api/v1/user_repos', { id: "1" }, {
          withCredentials: true
        });
        const repos = response.data;

        // Get usernames for each repository owner
        for (let repo of repos) {
          const userResponse = await axios.post('https://api.xvelper.ru/api/v1/user_details', { id: repo.OwnerID });
          repo.username = userResponse.data.username;
        }

        this.repositories = repos;
      } catch (error) {
        console.error('Error fetching repositories:', error);
      }
    },
    async deleteRepo(repoId, OwnerID) {
      try {
        const token = localStorage.getItem('token');
        console.log(token)
        const response = await axios.post('https://api.xvelper.ru/api/v1/delete_repo', {
          id: repoId,
          owner_id: OwnerID
        }, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });
        console.log(response);
        this.fetchRepositories(); // Обновить список репозиториев после удаления
      } catch (error) {
        console.error('Error deleting repository:', error);
      }
    },
    async createRepo() {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.post('https://api.xvelper.ru/api/v1/create_repo', {
          repo_name: this.newRepoName,
          is_public: this.isPublic,
        }, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });
        console.log('Repository created:', response.data);
        this.fetchRepositories(); // Обновить список репозиториев после создания
      } catch (error) {
        console.error('Error creating repository:', error);
      }
    },
    async viewCommits(repoId) {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.post('https://api.xvelper.ru/api/v1/get_commits', {
          repo_id: repoId,
        }, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });

        this.commitLogs = response.data;
        console.log(this.commitLogs);
      } catch (error) {
        console.error('Error fetching commits:', error);
      }
    },
    viewRepoContent(repoId) {
      this.$router.push({ name: 'UserRepoContent', params: { repoId } });
    }
  },
  created() {
    this.fetchRepositories();
  }
};
</script>

<style scoped>
.commits-head {
  display: flex;
  justify-content: space-between;
}

.btn_commit {
  margin-left: 5px;
}

.repoName {
  font-size: 26px;
  font-weight: 700;
}
.card {
  margin-bottom: 20px;
}
</style>