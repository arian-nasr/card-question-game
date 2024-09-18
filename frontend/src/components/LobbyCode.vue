<template>
    <form @submit.prevent="handleSubmit">
        <div class="mb-3">
            <label for="lobbycodeinput1" class="form-label">Lobby Code</label>
            <input type="text" class="form-control" id="lobbycodeinput1" v-model="lobbyCode">
            <label for="nameinput1" class="form-label">Your Name</label>
            <input type="text" class="form-control" id="nameinput1" v-model="name">
        </div>
        <button type="submit" class="btn btn-primary">Submit</button>
    </form>
</template>

<style>
div {
    color: white;
}
</style>
<script>
import axios from 'axios';

export default {
    data() {
        return {
            lobbyCode: '',
            name: ''
        };
    },
    methods: {
        async handleSubmit() {
            try {
                const formData = new FormData();
                formData.append('lobbyCode', this.lobbyCode);
                formData.append('name', this.name);

                const response = await axios.post('http://127.0.0.1:5000/join', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                console.log(response.data);
            } catch (error) {
                console.error('There was an error!', error);
            }
        },
        emitLobbyCode() {
            this.$emit('lobbyCode', this.lobbyCode);
        }
    }
};
</script>