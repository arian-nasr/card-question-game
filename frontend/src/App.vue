<template>
  <main>

    <Lobby v-if="showLobby" :lobbyCode="lobbyCode" :players="players" @start-game="startGame" />

    <div v-if="showLogin" class="showLogin">
      <form @submit.prevent="handleSubmit">
          <div class="mb-3">
              <label for="lobbycodeinput1" class="form-label">Lobby Code</label>
              <input type="text" class="form-control" id="lobbycodeinput1" v-model="lobbyCode">
              <label for="nameinput1" class="form-label">Your Name</label>
              <input type="text" class="form-control" id="nameinput1" v-model="name">
          </div>
          <button type="submit" class="btn btn-primary" ref="submitButton">Submit</button>
      </form>
    </div>

    <VotingPhaseNotify v-if="showVotingPhaseNotify" />

    <div v-if="showVotingCards" class="cardholder">
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

    <CommonCards v-if="showCommonCards" :commonQuestions="commonQuestions" />
    
  </main>
</template>
<style>
body {
  background-color: rgb(28, 31, 34);
}
.showLobby {
    color: white
}
.showLogin {
    color: white
}
.voting-phase-notify {
    text-align: center;
    color: white;
}
.cardholder {
    display: flex;
    justify-content: space-around;
    align-items: center;
    gap: 1rem;
}
.card-footer {
    display: flex;
    justify-content: space-around;
    gap: 1rem;
}
.card-body {
    text-align: center;
}
</style>
<script>
import axios from 'axios';

import Lobby from './components/Lobby.vue';
import VotingPhaseNotify from './components/VotingPhaseNotify.vue';
import CommonCards from './components/CommonCards.vue';

const API_URL = 'http://192.168.2.158:5000';

export default {
  components: {
    Lobby,
    VotingPhaseNotify,
    CommonCards
  },
  data() {
    return {
      lobbyCode: '',
      name: '',
      showLogin: true,
      showLobby: false,
      showVotingPhaseNotify: false,
      showVotingCards: false,
      questions: [], // New property to store questions
      approvedQuestions: [], // New property to store approved questions
      rejectedQuestions: [], // New property to store rejected questions
      showCommonCards: false,
      commonQuestions: [], // New property to store common questions
      players: [] // New property to store players
    }
  },
  methods: {
    async voteOnCard(index, accept) {
      if (accept) {
        // Check if the question is not already in the approved list
        if (!this.approvedQuestions.includes(this.questions[index])) {
          this.approvedQuestions.push(this.questions[index]); // Add the question to the approved list if not already present
        }
      }
      if (!accept) {
        // Check if the question is not already in the rejected list
        if (!this.rejectedQuestions.includes(this.questions[index])) {
          this.rejectedQuestions.push(this.questions[index]); // Add the question to the rejected list if not already present
        }
      }
      // Check if all questions have been voted on
      if (this.approvedQuestions.length + this.rejectedQuestions.length === this.questions.length) {
        // Call a method to send the approved questions to the server
        await this.sendApprovedQuestions(); 
      }
    },

    async startGame() {
      // Implement this method to start the game after all players have joined the lobby
      try {
        const formData = new FormData();
        formData.append('lobbyCode', this.lobbyCode);
        const response = await axios.post(`${API_URL}/startVoting`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
      } catch (error) {
        console.error('Error sending approved questions:', error);
      }
    },
    async sendApprovedQuestions() {
      try {
        const formData = new FormData();
        formData.append('lobbyCode', this.lobbyCode);
        formData.append('name', this.name);
        formData.append('approvedQuestions', JSON.stringify(this.approvedQuestions));
        const response = await axios.post(`${API_URL}/approveQuestions`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        this.showVotingCards = false;
        // Call the method to poll game status at this point
        await this.pollGameStatus();
      } catch (error) {
        console.error('Error sending approved questions:', error);
      }
    },

    async pollPlayers() {
      try {
      const formData = new FormData();
      formData.append('lobbyCode', this.lobbyCode);

      const poll = async () => {
        try {
        const response = await axios.post(`${API_URL}/checkPlayers`, formData, {
          headers: {
          'Content-Type': 'multipart/form-data'
          }
        });
        
        if (response.status === 200 && response.data.status === 'waiting') {
          this.players = response.data.players; // Assuming the response contains a 'players' array
          setTimeout(poll, 1000); // Poll again after 1 second
        } else
        if (response.status === 200 && response.data.status === 'voting') {
          this.showLobby = false;
          this.showVotingPhaseNotify = true;
          await this.getVotingQuestions();
          this.showVotingCards = true;
        } else {
          setTimeout(poll, 1000); // Poll again after 1 second
        }
        } catch (error) {
        console.error('Error polling players:', error);
        setTimeout(poll, 1000); // Poll again after 1 second in case of error
        }
      };

      poll(); // Start the polling loop
      } catch (error) {
      console.error('Error initiating players poll:', error);
      }
    },

    // This method will be called when a game is joined
    async handleSubmit() {
      try {
      const formData = new FormData();
      formData.append('lobbyCode', this.lobbyCode);
      formData.append('name', this.name);

      const response = await axios.post(`${API_URL}/joinGame`, formData, {
        headers: {
        'Content-Type': 'multipart/form-data'
        }
      });

      if (response.status === 200) {
        this.showLogin = false;
        this.showLobby = true;
        await this.pollPlayers();
      } else {
        console.error('Failed to join the game:', response.data);
        this.$refs.submitButton.innerText = 'Failed to join the game';
      }
      } catch (error) {
        console.error('Error joining the game:', error);
        this.$refs.submitButton.innerText = 'Error joining the game';
      }
    },

    // This method will be called every few seconds to poll the game status
    async pollGameStatus() {
      try {
      const formData = new FormData();
      formData.append('lobbyCode', this.lobbyCode);

      const poll = async () => {
        try {
        const response = await axios.post(`${API_URL}/checkStatus`, formData, {
          headers: {
          'Content-Type': 'multipart/form-data'
          }
        });
        
        if (response.status === 200 && response.data === 'ingame') {
          await this.getCommonQuestions(); // Call the method to get common questions
        } else {
          setTimeout(poll, 1000); // Poll again after 1 second
        }
        } catch (error) {
        console.error('Error polling game status:', error);
        setTimeout(poll, 1000); // Poll again after 1 second in case of error
        }
      };

      poll(); // Start the polling loop
      } catch (error) {
      console.error('Error initiating game status poll:', error);
      }
    },

    async getVotingQuestions() {
      try {
        const formData = new FormData();
        formData.append('lobbyCode', this.lobbyCode);
        // We probably don't need to send the name here since both users in the same lobby will get the same questions
        //formData.append('name', this.name);

        const response = await axios.post(`${API_URL}/getQuestions`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });

        if (response.status === 200) {
          this.questions = response.data; // Assuming the response contains a 'questions' array
        } else {
          console.error('Failed to get voting questions:', response.data);
        }
      } catch (error) {
        console.error('Error getting voting questions:', error);
      }
    },
    async getCommonQuestions() {
      try {
      const formData = new FormData();
      formData.append('lobbyCode', this.lobbyCode);

      const response = await axios.post(`${API_URL}/getCommonQuestions`, formData, {
        headers: {
        'Content-Type': 'multipart/form-data'
        }
      });

      // Assuming the response contains a 'commonQuestions' array, add the response to the commonQuestions array
      if (response.status === 200) {
        this.commonQuestions = response.data;
        this.showVotingPhaseNotify = false;
        this.showCommonCards = true;
      } else {
        console.error('Failed to get common questions:', response.data);
      }
      } catch (error) {
      console.error('Error getting common questions:', error);
      }
    }
    },
};
</script>