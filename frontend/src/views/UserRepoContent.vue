<template>
  <div class="container mt-5">
    <table class="table table-hover">
      <thead>
        <tr>
          <th>Repository: {{ repoName }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="file in repoContents" :key="file">
          <td>{{ file }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'UserRepoContent',
  data() {
    return {
      repoContents: [],
      repoName: '',
    };
  },
  methods: {
    async fetchRepoContents(repoId) {
      try {
        const token = localStorage.getItem('token');
        console.log(repoId)
        const response = await axios.post('https://api.xvelper.ru/api/v1/view_contents', {
          repo_id: repoId,
        }, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });
        this.repoContents = response.data.contents;
        this.repoName = this.$route.params.repoName; // Получаем имя репозитория из параметров маршрута
      } catch (error) {
        console.error('Error fetching repository contents:', error);
      }
    }
  },
  created() {
    const repoId = this.$route.params.repoId;
    this.fetchRepoContents(repoId);
  }
};
</script>

<style scoped>
.table-hover tbody tr:hover {
  background-color: #f5f5f5;
}
</style>
