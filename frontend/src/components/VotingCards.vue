<template>
    <div class="cardholder">
        <!-- Dynamically load question cards -->
        <div v-for="(question, index) in questions" :key="index" class="card">
            <div class="card-body">
                {{ question }}
            </div>
            <div class="card-footer">
                <a href="#" class="card-link" @click.prevent="voteOnCard(index, true)">Accept</a>
                <a href="#" class="card-link" @click.prevent="voteOnCard(index, false)">Decline</a>
            </div>
        </div>
    </div>
</template>
<script>
export default {
  props: {
    questions: Array
  },
  methods: {
    voteOnCard(index, vote) {
      this.$emit('vote-on-card', index, vote);
    }
  }
}
</script>