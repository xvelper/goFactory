<template>
  <div class="container mt-5">
    <h2 class="text-center">Добро пожаловать в GitFactory</h2>
    <div class="row justify-content-center mt-4">
      <div class="col-md-6 mb-4 d-flex justify-content-center" v-for="(info, index) in infoBlocks" :key="index">
        <div class="card h-100" style="width: 80%;">
          <div class="card-body">
            <h5 class="card-title text-center">{{ info.title }}</h5>
            <p class="card-text">{{ info.text }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomePage',
  data() {
    return {
      infoBlocks: [
        {
          title: 'Простой в установке',
          text: 'Просто запустите исполняемый файл для вашей платформы. Используйте GoFactory с Docker.'
        },
        {
          title: 'Кроссплатформенный',
          text: 'GoFactory работает на любой операционной системе, которая может компилировать Go: Windows, macOS, Linux, ARM и т. д. Выбирайте, что вам больше нравится!'
        },
        {
          title: 'Легковесный',
          text: 'GoFactory имеет низкие системные требования и может работать на недорогом Raspberry Pi. Экономьте энергию вашей машины!'
        },
        {
          title: 'Открытый исходный код',
          text: 'Всё это на GitHub! Чтобы получить доступ обратитесь в ТГ: @idtym!'
        }
      ]
    };
  }
};
</script>

<style scoped>
.text-center {
  text-align: center;
}
</style>
