<template>
  <div class="container mt-5">

    <div class="d-flex">
      <div class="col-md-4 text-center">
        <img :src="avatarUrl" class="rounded-circle avatar mb-3" alt="User Avatar" />

        <p>{{ username }}</p>
      </div>
      <div class="row">

        <div class="p-2 flex-grow-1"><h4>Данные пользователя</h4></div>
        <div class="row">
          <div class="d-flex">
          <div class="row">
            <h5>Имя</h5>
            <div class="d-flex">
              <p>{{ username }}</p>
              <button type="button" class="btn btn-secondary btn-sm">Ред.</button>
            </div>
          </div>
          <div class="row">
            <h5>Фамилия</h5>
            <div class="d-flex">
              <div class="input-group mb-3">
                <input type="text" class="form-control" aria-label="Фамилия" aria-describedby="inputGroup-sizing-default">
              </div>
              <button type="button" class="btn btn-secondary btn-sm">Ред.</button>
            </div>
          </div>
        </div>
        <div class="d-flex">
          <div class="row">
            <h5>Имя</h5>
            <div class="d-flex">
              <p>{{ username }}</p>
              <button type="button" class="btn btn-secondary btn-sm">Ред.</button>
            </div>
          </div>
          <div class="row">
            <h5>Фамилия</h5>
            <div class="d-flex">
              <p>{{ username }}</p>
              <button type="button" class="btn btn-secondary btn-sm">Ред.</button>
            </div>
          </div>
        </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'UserProfile',
  data() {
    return {
      avatarUrl: 'https://via.placeholder.com/150',
      fullName: '',
      username: '',
      popularRepos: [],
    };
  },
  methods: {
    async fetchUserDetails() {
      try {
        const token = localStorage.getItem('token');
        console.log('Fetching user details with token:', token);
        const response = await axios.get('http://localhost:8000/api/v1/user_details', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        console.log('User details response:', response.data);
        this.fullName = response.data.username;
        this.username = response.data.username;
      } catch (error) {
        console.error('Error fetching user details', error);
      }
    },
    async fetchPopularRepos() {
      try {
        const token = localStorage.getItem('token');
        console.log('Fetching popular repos with token:', token);
        const response = await axios.get('http://localhost:8000/api/v1/user_repos', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        console.log('Popular repos response:', response.data);
        this.popularRepos = response.data;
      } catch (error) {
        console.error('Error fetching popular repositories', error);
      }
    },
    viewCommits(uuid) {
      console.log(`Viewing commits for repository with UUID: ${uuid}`);
      // Implement logic to view commits
    },
    async deleteRepository(uuid) {
      try {
        const token = localStorage.getItem('token');
        console.log('Deleting repository with token:', token);
        await axios.delete(`http://localhost:8000/api/v1/delete_repo`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
          data: {
            uuid: uuid,
          },
        });
        console.log('Repository deleted successfully');
        this.fetchPopularRepos(); // Refresh the list of repositories
      } catch (error) {
        console.error('Error deleting repository', error);
      }
    },
  },
  created() {
    this.fetchUserDetails();
    this.fetchPopularRepos();
  },
};
</script>

<style scoped>
.avatar {
  width: 150px;
  height: 150px;
  border: 3px solid #343a40;
}

.repo-card {
  background-color: #fff;
  border: 1px solid #ced4da;
}

.language-color {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 5px;
}

.language-color.JavaScript {
  background-color: #f1e05a;
}

.language-color.HTML {
  background-color: #e34c26;
}

.language-color.Python {
  background-color: #3572A5;
}

.language-color.PHP {
  background-color: #4F5D95;
}
</style>