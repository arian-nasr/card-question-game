from flask import Flask, request, json, jsonify
from flask_cors import CORS
from pymongo import MongoClient
import random


# Connect to MongoDB
client = MongoClient('mongodb://localhost:27017/')
db = client["game-server-1"]
collection = db["activegames"]


app = Flask(__name__)
CORS(app)


# The following POST api will be used to join a game given a lobbyCode and a name, if the lobbyCode is not in the collection, return 'Invalid lobby code'
# If user is already in the game, do not add them again
@app.route('/joinGame', methods=['POST'])
def joinGame():
    lobbyCode = request.form['lobbyCode']
    name = request.form['name']
    if collection.find_one({"lobbyCode": lobbyCode}):
        players = collection.find_one({"lobbyCode": lobbyCode})['players']
        if name in players:
            return 'User already in game'
        else:
            players.append(name)
            collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"players": players}})
            return 'Success'
    else:
        return 'Lobby not found', 400
    

# The following POST api will be used to create a game given a lobbyCode and a name, if the lobbyCode is already in the collection, return 'Lobby code already exists'
@app.route('/createGame', methods=['POST'])
def createGame():
    lobbyCode = request.form['lobbyCode']
    name = request.form['name']
    if collection.find_one({"lobbyCode": lobbyCode}):
        return 'Lobby code already exists'
    else:
        # Set the game status to 'voting' here, this might be changed later
        collection.insert_one({"lobbyCode": lobbyCode, "players": [name], "status": 'waiting'})
        return 'Success'
    

# The following function will be used to get 3 random questions from the questions.txt file
# Repeat questions are not being checked for, so it is possible to get the same question multiple times
# This should be fixed in a future update
def get_random_questions():
    with open('./questions.txt', 'r') as file:
        lines = file.readlines()
        lines = [line.strip() for line in lines]  # Remove newline characters
        return random.sample(lines, 3)
    

# The following POST api will be used to get 3 random questions, questions should be saved in the collection
@app.route('/getQuestions', methods=['POST'])
def getQuestions():
    lobbyCode = request.form['lobbyCode']
    game = collection.find_one({"lobbyCode": lobbyCode})
    if game:
        if 'questions' in game:
            return game['questions']
        else:
            questions = get_random_questions()
            collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"questions": questions}})
            return questions
    else:
        return 'Invalid lobby code'
    
    
# The following function will compare the approved questions and return the common questions between all players
def getCommonApproved(lobbyCode):
    game = collection.find_one({"lobbyCode": lobbyCode})
    approvedQuestions = game['approvedQuestions']
    # create a set of the first player's approved questions
    commonQuestions = set(approvedQuestions[game['players'][0]])
    # iterate through the rest of the players and find the common questions
    for player in game['players'][1:]:
        commonQuestions = commonQuestions.intersection(approvedQuestions[player])
    collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"commonQuestions": list(commonQuestions)}})

    # voting phase is now over and common questions have been calculated
    # add game status to collection, set to 'ingame'
    collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"status": 'ingame'}})


# The following POST api will be used to check the status of the game
# This endpoint will be used to poll the server to check if the game should advance to the next phases
@app.route('/checkStatus', methods=['POST'])
def checkStatus():
    lobbyCode = request.form['lobbyCode']
    game = collection.find_one({"lobbyCode": lobbyCode})
    if game:
        return game['status']
    else:
        return 'Invalid lobby code'
    

# The following POST api will be used to check the players in the game
@app.route('/checkPlayers', methods=['POST'])
def checkPlayers():
    lobbyCode = request.form['lobbyCode']
    game = collection.find_one({"lobbyCode": lobbyCode})
    if game:
        return jsonify(
            players=game['players'],
            status=game['status']
        )
    else:
        return 'Invalid lobby code'
    

# The following POST api will be used to set the status of the game to voting
@app.route('/startVoting', methods=['POST'])
def startVoting():
    lobbyCode = request.form['lobbyCode']
    collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"status": 'voting'}})
    return 'Success'


# Users will approve questions by sending a POST request to this api with their approved questions
# This function will add the approved questions to the collection for each user so they can be compared later
@app.route('/approveQuestions', methods=['POST'])
def approveQuestions():
    lobbyCode = request.form['lobbyCode']
    name = request.form['name']
    questions = request.form['approvedQuestions']
    questions = json.loads(questions)  # Parse the JSON string into a Python list
    game = collection.find_one({"lobbyCode": lobbyCode})
    if game:
        if 'approvedQuestions' in game:
            approvedQuestions = game['approvedQuestions']
            approvedQuestions[name] = questions
            collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"approvedQuestions": approvedQuestions}})
        else:
            collection.update_one({"lobbyCode": lobbyCode}, {"$set": {"approvedQuestions": {name: questions}}})
        
        # Check if all players have approved questions
        game = collection.find_one({"lobbyCode": lobbyCode})
        players = game['players']
        approvedQuestions = game['approvedQuestions']
        if len(players) > 1:
            if all(player in approvedQuestions for player in players):
                getCommonApproved(lobbyCode)
        
        return 'Success'
    else:
        return 'Invalid lobby code'


# The following POST api will be used to get the common questions between all players
@app.route('/getCommonQuestions', methods=['POST'])
def getCommonQuestions():
    lobbyCode = request.form['lobbyCode']
    game = collection.find_one({"lobbyCode": lobbyCode})
    if game:
        return game['commonQuestions']
    else:
        return 'Invalid lobby code'