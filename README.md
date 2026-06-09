# Introduction
    
&nbsp;&nbsp;&nbsp;&nbsp;This project aims to unify data from the IT level(ERP) with the OT level(Ignition/factory floor) through a common middleware(the UNS).
    
# Data flow

```text
┌─────────────────────────┐                             
│                         │                             
│OT Layer(SCADA Ignition) ◄──┐                          
│                         │  │                          
└─────────────────────────┘  │                          
                             │                          
                             │                          
            ┌────────────────┘                          
            │                                           
┌───────────▼─────────────┐                             
│                         │                             
│ MQTT Mosquitto broker   ◄──┐                          
│                         │  │                          
└─────────────────────────┘  │                          
                             │                          
            ┌────────────────┘                          
            │                                           
            │                                           
┌───────────▼─────────────┐      ┌─────────────────────┐
│                         │      │                     │
│       UNS in Go         ◄──────►      ERP in Java    │
│                         │      │                     │
└───────────▲─────────────┘      └─────────────────────┘
            │                                           
            │                                           
            │                                           
┌───────────▼─────────────┐                             
│                         │                             
│   Postgresql Database   │                             
│                         │                             
└─────────────────────────┘                             
```
&nbsp;&nbsp;&nbsp;&nbsp;Diagram made with https://asciiflow.com/#/
    
# Considerations
    
&nbsp;&nbsp;&nbsp;&nbsp;Using this architecture, the UNS made with Go acts like a server in a client-server application. The app gets MQTT messages from the OT level and fetches the corresponding ERP order data for each conveyor(the OT level uses a PLC simulator) then pushes the new data to a Postgresql DB with a timestamp for when it was pushed and timestamps from the OT level and ERP. The timestamp format used is RFC3339.\
&nbsp;&nbsp;&nbsp;&nbsp;The ERP can also make order requests and send them to the UNS through HTTP and then the UNS forwards it to the MQTT broker.\When the OT finishes processing the order it sends a message to the broker and the UNS picks it up and forwards it to the ERP so it can update it's web UI.\
&nbsp;&nbsp;&nbsp;&nbsp;Docker was also used because it simplified deployment and iteration testing. Before using docker, I had to always start mosquitto, postgresql, pgadmin and run the uns code, this took a lot of time.\
&nbsp;&nbsp;&nbsp;&nbsp;For testing remotely we used Hamachi to create a VPN.\
&nbsp;&nbsp;&nbsp;&nbsp;The app could also be deployed to a VPS but with some security changes in mind, currently the mosquitto broker is configured to listen on all interfaces and allow anonymous connections.
                                                             
# Usage
    
git clone <this_project_url>\
cd <this_project_url>\
touch .env(edit .env file with required env variables)\
docker compose up --build\
-> localhost:5050 to access db\
-> development credentials: admin@uns.com admin\
     
# Requirements
    
docker desktop && docker\
postgresql && pgadmin\
