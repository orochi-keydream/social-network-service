#!/usr/bin/env tarantool
-- Configure database
box.cfg {
    listen = 3301
}

box.once("bootstrap", function()
    box.schema.sequence.create('message_id_seq', { start = 1 })

    box.schema.space.create('messages')

    box.space.messages:format(
        {
            { name = 'message_id', type = 'unsigned' },
            { name = 'from_user_id', type = 'string' },
            { name = 'to_user_id', type = 'string' },
            { name = 'sent_at', type = 'string' },
            { name = 'text', type = 'string' }
        })

    box.space.messages:create_index('primary',
        { type = 'TREE', parts = {1, 'unsigned'}, sequence = 'message_id_seq' })

    box.space.messages:create_index('from_user_id',
        { type = 'TREE', parts = {2, 'string'}, unique = false })

    box.space.messages:create_index('to_user_id',
        { type = 'TREE', parts = {3, 'string'}, unique = false })
end)

function get_messages(currentUser, targetUser)
    params = {{[':a'] = currentUser}, {[':b'] = targetUser}}
    return box.execute(
        [[select * from messages where (from_user_id = :a and to_user_id = :b) or (from_user_id = :b and to_user_id = :a) order by message_id desc;]],
        params)
end

function add_message(currentUser, targetUser, sentAt, text)
    return box.space.messages:insert({nil, currentUser, targetUser, sentAt, text})
end