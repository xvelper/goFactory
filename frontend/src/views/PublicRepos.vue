<template>
  <div class="container mt-5">
    <h2 class="text-center">Открытые репозитории</h2>
    <div class="row">
      <div v-for="repo in repos" :key="repo.ID" class="col-md-4 my-3">
        <div class="card">
          <div class="card-body">
            <h5 class="card-title">{{ repo.Name }}</h5>
            <p class="card-text">URL: http://localhost:8000/{{ repo.Username }}/{{ repo.Name }}</p>
            <p class="card-text">Owner: {{ repo.Username }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'PublicRepos',
  data() {
    return {
      repos: []
    };
  },
  methods: {
    async fetchRepos() {
      try {
        const response = await axios.get('https://api.xvelper.ru/api/v1/public_repos');
        const repos = response.data;

        for (let repo of repos) {
          const userResponse = await axios.post('https://api.xvelper.ru/api/v1/user_details', { id: repo.OwnerID });
          repo.Username = userResponse.data.username;
        }

        this.repos = repos;
      } catch (error) {
        console.error('Error fetching repositories:', error);
      }
    }
  },
  created() {
    this.fetchRepos();
  }
};
</script>

<style scoped>
.card {
  width: 100%;
}
.text-center {
  text-align: center;
}
</style>
