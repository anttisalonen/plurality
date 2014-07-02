#include <fstream>
#include <iostream>
#include <vector>
#include <boost/shared_ptr.hpp>

#include <jsoncpp/json/json.h>


class Component {
	public:
		Component(const std::string& name) : mName(name) { }
		~Component() { }
		virtual void Start() { }

	protected:
		std::string mName;
};

typedef boost::shared_ptr<Component> ComponentPtr;

class GameObject {
	public:
		GameObject(const std::string& name) : mName(name) { }
		void addComponent(ComponentPtr c) { mComponents.push_back(c); }
		std::vector<ComponentPtr> getComponents() { return mComponents; }

	private:
		std::string mName;
		std::vector<ComponentPtr> mComponents;
};

class HelloComponent : public Component {
	public:
		HelloComponent();
		virtual void Start() override;
};

HelloComponent::HelloComponent()
	: Component("HelloComponent")
{
}

void HelloComponent::Start()
{
	std::cout << "Hello world!\n";
}

void runGame(const Json::Value& root)
{
	
	std::vector<GameObject> Objects;
	for(auto& jo : root["objects"]) {
		GameObject obj(jo["name"].asString());
		for(auto& jcomp : jo["components"]) {
			const std::string& type = jcomp["type"].asString();
			if(type == "HelloComponent") 
				obj.addComponent(ComponentPtr(new HelloComponent()));
			else
				std::cerr << "Invalid component type " << type << "!\n";
		}
		Objects.push_back(obj);
	}

	for(auto& obj : Objects) {
		for(auto& comp : obj.getComponents())
			comp->Start();
	}

	Objects.clear();
}

int main(int argc, char** argv)
{
	if(argc != 2) {
		std::cerr << "Usage: " << argv[0] << " <game JSON file>\n";
		exit(1);
	}
	std::string jsonFilename = argv[1];

	Json::Reader reader;
	Json::Value root;

	std::ifstream input(jsonFilename, std::ifstream::binary);
	bool parsingSuccessful = reader.parse(input, root, false);
	if (!parsingSuccessful) {
		throw std::runtime_error(reader.getFormatedErrorMessages());
	}

	runGame(root);

	return 0;
}
