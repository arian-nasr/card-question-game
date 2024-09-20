<template>
    <div class="cardholder">
        <!-- Dynamically load common cards -->
        <div v-for="(question, index) in commonQuestions" :key="index" class="card">
            <div class="card-body">
                {{ question }}
            </div>
        </div>
    </div>
</template>
<script>
  export default {
    props: {
      commonQuestions: Array
    },
  };
  </script>