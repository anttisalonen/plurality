#include <fstream>
#include <iostream>
#include <vector>
#include <map>

#include <boost/shared_ptr.hpp>
#include <boost/variant.hpp>

#include <jsoncpp/json/json.h>


class Component {
	public:
		Component(const std::string& name) : mName(name) { }
		~Component() { }
		virtual void Start() { }
		void addValue(std::string name, const std::string& value);
		void addValue(std::string name, int value);

	protected:
		std::string mName;
		std::map<std::string, boost::variant<std::string, int>> mValues;
};

void Component::addValue(std::string name, const std::string& value)
{
	mValues[name] = value;
}

void Component::addValue(std::string name, int value)
{
	mValues[name] = value;
}

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
	const std::string& g = boost::get<std::string>(mValues["greetee"]);
	int num = boost::get<int>(mValues["number of greets"]);
	for(int i = 0; i < num; i++)
		std::cout << "Hello " << g << "!\n";
}

void runGame(const Json::Value& root)
{
	
	std::vector<GameObject> Objects;
	for(auto& jo : root["objects"]) {
		GameObject obj(jo["name"].asString());
		for(auto& jcomp : jo["components"]) {
			const std::string& type = jcomp["type"].asString();
			ComponentPtr comp;
			if(type == "HelloComponent") 
				comp = ComponentPtr(new HelloComponent());
			else
				std::cerr << "Invalid component type " << type << "!\n";

			if(comp) {
				obj.addComponent(comp);
				auto jvalnames = jcomp["values"].getMemberNames();
				for(auto& jvalname : jvalnames) {
					const std::string& valname = jvalname;
					const Json::Value& value = jcomp["values"][valname];
					if(value.isString())
						comp->addValue(valname, value.asString());
					else if(value.isIntegral())
						comp->addValue(valname, value.asInt());
					else
						std::cerr << "Invalid value type!\n";
				}
			}
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
