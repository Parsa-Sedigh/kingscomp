-- check queue size and match users
local function matchUsers(queueKey, pubSubChannel, minUsers, minScore, lobbyId, userId, userScore)
    local users = redis.call('ZRANGEBYSCORE', queueKey, minScore, '+inf', 'LIMIT', 0, minUsers)

    if #users >= minUsers then
        for i, v in ipairs(users) do
            users[i] = tonumber(v)
        end

        table.insert(users, userId)

        -- remove these users from sorted set
        redis.call('ZREM', queueKey, unpack(users))

        local lobby = {
            id = lobbyId,
            participants = users,
            created_at = userScore,
            state = 'started'
        }
        local lobbyJson = cjson.encode(lobby)

        for i, v in ipairs(users) do
            redis.call('JSON.MSET')
        end

        -- notify the matched users via pub/sub channel
        redis.call('PUBLISH', pubSubChannel, lobbyId .. ':' .. table.concat(users, ','))

        return {true, lobbyId, users} -- matching succeeded
    end

    -- add the current user to the queue since not enough users are present
    redis.call('ZADD', queueKey, userScore, userId)

    return {false} -- not enough users for matching
end

-- keys and arguments
local queueKey = KEYS[1]
local pubSubChannel = KEYS[2]
local minUsers = tonumber(ARGV[1])
local minScore = tonumber(ARGV[2])
local lobbyId = ARGV[3]
local userId = tonumber(ARGV[4])
local userScore = tonumber(ARGV[5])

return matchUsers(queueKey, pubSubChannel, minUsers, minScore, lobbyId, userId, userScore)