    # Introduction
    
    This project aims to unify data from the IT level(ERP) with the OT level(Ignition/factory floor) through a common middleware(the UNS).
    
    # Data flow
    
                                                             
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
│       UNS in Go         ◄──────►     ERP in Java     │ 
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
                                                                 
    Diagram made with https://asciiflow.com/#/
    
    # Considerations
    
    Using this architecture, the UNS made with Go acts like a server in a client-server application. The app gets MQTT messages from the OT level and fetches the corresponding ERP order data for each conveyor(the OT level uses a PLC simulator) then pushes the new data to a Postgresql DB with a timestamp for when it was pushed and timestamps from the OT level and ERP. The timestamp format used is RFC3339.
                                                             
    # Usage
    
    git clone <this_project_url>
    cd <this_project_url>
    touch .env(edit .env file with required env variables)
    docker compose up --build
     -> localhost:5050 to access db
     -> development credentials: admin@uns.com admin
     
    # Requirements
    
    docker desktop && docker
    postgresql && pgadmin
