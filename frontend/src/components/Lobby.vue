<template>
    <div class="showLobby">
      <h1>Lobby Code: {{ lobbyCode }}</h1>
      <button class="btn btn-primary" @click="startGame">Start Game</button>
      <p>Players:</p>
      <p v-for="(player, index) in players" :key="index">
        {{ player }}
      </p>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      lobbyCode: String,
      players: Array,
    },
    methods: {
      startGame() {
        this.$emit('start-game');
      },
    },
  };
  </script>