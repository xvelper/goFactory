import { createRouter, createWebHistory } from 'vue-router';
import HomePage from '../views/HomePage.vue';
import UserRegister from '../views/UserRegister.vue';
import UserLogin from '../views/UserLogin.vue';
import UserRepos from '../views/UserRepos.vue';
import PublicRepos from '../views/PublicRepos.vue';
import UserProfile from '../views/UserProfile.vue';
import MyRepositories from '@/views/MyRepositories.vue';

const routes = [
  { path: '/', name: 'HomePage', component: HomePage },
  { path: '/register', name: 'UserRegister', component: UserRegister },
  { path: '/login', name: 'UserLogin', component: UserLogin },
  { path: '/user-repos', name: 'UserRepos', component: UserRepos },
  { path: '/public-repos', name: 'PublicRepos', component: PublicRepos },
  { path: '/profile', name: 'UserProfile', component: UserProfile },
  { path: '/my-repositories', name: 'MyRepositories', component: MyRepositories },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
