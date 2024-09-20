<template>
    <div class="showLogin">
        <form @submit.prevent="handleSubmit">
            <div class="mb-3">
                <label for="lobbycodeinput1" class="form-label">Lobby Code</label>
                <input type="text" class="form-control" id="lobbycodeinput1" v-model="localLobbyCode">
                <label for="nameinput1" class="form-label">Your Name</label>
                <input type="text" class="form-control" id="nameinput1" v-model="localName">
            </div>
            <button type="submit" class="btn btn-primary" ref="submitButton">{{ submitButtonText }}</button>
        </form>
    </div>
</template>
<script>
export default {
  props: {
    lobbyCode: String,
    name: String,
    submitButtonText: String
  },
  data() {
    return {
      localLobbyCode: this.lobbyCode,
      localName: this.name
    };
  },
  methods: {
    handleSubmit() {
      this.$emit('update:lobbyCode', this.localLobbyCode);
      this.$emit('update:name', this.localName);
      this.$emit('submit');
    }
  }
}
</script>