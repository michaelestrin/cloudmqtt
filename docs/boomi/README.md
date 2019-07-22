# michaelestrin/cloudmqtt/docs/boomi -- Connecting to Dell Boomi AtomSphere

I was able to connect the service to Dell Boomi AtomSphere using an intermediary MQTT broker.  

## CloudMQTT Instance Configuration

First, I created a free account and a new _Cute Cat_ instance on [CloudMQTT](https://cloudmqtt.com):

![CloudMQTT Instance](images/cloudmqtt.PNG)

## Boomi Process, Connector, and Operation Configuration

Next, within the _Build_ context in Boomi, I selected the __New__ option:

![New](images/new.PNG)

In the follow-on dialog, I selected _Connection_ for _Component Type_, _MQTT Connector_ for _Connector_ and 
    named the component "EdgeX" before selecting the __Create__ option:

![New Connection Dialog](images/newconnectiondialog.PNG)

I then entered the credentials and connection information from the CloudMQTT instance and selected the 
    __Save and Close__ option:

![New Connection Settings](images/newconnectionsettings.PNG)

Next, 

![New](images/new.PNG)

In the follow-on dialog, I named the component and selected the __Create__ option:

![Create Process](images/createprocess.PNG)

I then clicked on the __Make the recommended changes for me__ link:

![Create Recommendations](images/createrecommendations.PNG)

I then named the start shape "EdgeX Inbound", selected _MQTT Connector_ for _Connector_, _Listen_ for _Action_, 
    _EdgeX_ for _Connection_, and clicked the green plus option in the _Operation_ field:
    
![Create Settings](images/createsettings.PNG)

I named the MQTT Operation "Listen", set _Topic_ to _topic/events_, and seleced the __Save and Close__ option: 

![Listen Operation](images/listenoperation.PNG)

This took me back to the Start Shape dialog where I selected the __OK__ option:

![Create Settings](images/createsettings2.PNG)

This took me back to the EdgeX process screen.  Under _Search Shapes_, I selected _Connect_ found the 
    _MQTT Connector_ option, and dragged and dropped it onto the process screen.  This brought up a Connector Shape 
    Dialog where I set the Display Name to "EdgeX Outbound", _Send_ for _Action_, _EdgeX_ for _Connection_, and 
    and clicked the green plus option in the _Operation_ field:

![Create Settings](images/createsettings3.PNG)

I named the MQTT Operation "Send", set _Topic_ to _topic/commands_, and seleced the __Save and Close__ option: 

![Send Operation](images/sendoperation.PNG)

This took me back to the Start Shape dialog where I selected the __OK__ option:

![Create Settings](images/createsettings4.PNG)

This took me back to the EdgeX process screen:

![Create Process](images/createprocess2.PNG)

Under _Search Shapes_, I selected _Logic_ found the _Stop_ option, and dragged and dropped it onto the process 
    screen. This brought up a Stop Shape Dialog where I set the Display Name to "Stop", and selected the __OK__ option:

![Stop Operation](images/stopoperation.PNG)

This took me back to the EdgeX process screen:

![Create Process](images/createprocess3.PNG)

I then connected the components by dragging the red arrows and selected the __Save and Close__ option:

![Create Process](images/createprocess4.PNG)

## Boomi Atom and Environment Configuration

In the _Manage_ context in Boomi, I selected the __Add Your First Environment__ option:

![First Environment](images/firstenvironment.PNG)

I then entered "EdgeX Test Environment" for the new environment's name, _Test_ for _Environment Classification_, and 
    selected the __Save__ option:

![First Environment](images/firstenvironmentdialog.PNG)

This brought me back to the _Environments_ screen where I selected the __New Atom__ option:

![Environments](images/environments.PNG)

I then downloaded and installed the Atom installer:

![Atom Installer](images/atominstaller.PNG)

When I subsequently refreshed the _Environments_ page, my newly created Atom appeared:

![Atom Installed](images/atominstalled.PNG)

I then selected the newly created Atom and -- in the _Attach to Environment_ dropdown under _Atom Controls_ -- the 
    _EdgeX Test Environment_ I created earlier:

![Atom Attach](images/atomattach.PNG)

The following screen displayed after the Atom had been successfully attached to the environment:

![Atom Attached](images/atomattached.PNG)

## Deploying the Process

Within the _Deploy_ context in Boomi, I checked the _EdgeX_ process and selected the __Save and Deploy__ option:

![Deploy](images/deploy.PNG)

This brought up a dialog where I selected the __OK__ option:

![Deploy](images/deploy2.PNG)

The process was subsequently deployed to the Atom.  I verified the process had been successfully deployed to the 
    Atom by going to the _Manage > Atom Management_ context:

![Deploy](images/deployed.PNG)

I then selected the Atom, selected the _Listeners_ option, and verified the _MQTT Connector > EdgeX_ listener had a 
    green circle next to it:
    
![Deploy](images/deployed2.PNG)

## configuration.toml Settings

I then updated my configuration.toml settings as follows:

![Settings](images/settings.PNG)
