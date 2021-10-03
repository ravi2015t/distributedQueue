# distributedQueue
A distributed queue like Kafka

# Features (in progress)

1. Easy to configure out of the box to not lose data. 
2. Distributed, with replication out of the box. 


# Design (in progress)

1. Data is split into chunks and stored as files on disk. 
2. Client readers read the chunks by specifying the offset. 
3. Readers explicitly acknowledge the data that was read. 
4. Every message is separated by \n End of line. 



